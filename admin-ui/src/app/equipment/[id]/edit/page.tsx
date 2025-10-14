"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { equipmentApi } from "@/lib/api/equipment";
import type { Equipment } from "@/types";
import { ArrowLeft, Loader2 } from "lucide-react";

const schema = z.object({
  equipment_name: z.string().min(2, "Required"),
  manufacturer_name: z.string().min(2, "Required"),
  serial_number: z.string().min(1, "Required"),
  model_number: z.string().optional(),
  category: z.string().optional(),
  customer_name: z.string().min(2, "Required"),
  installation_location: z.string().optional(),
  installation_date: z.string().optional(),
  purchase_date: z.string().optional(),
  warranty_expiry: z.string().optional(),
  notes: z.string().optional(),
});

type FormValues = z.infer<typeof schema>;

export default function EditEquipmentPage({ params }: { params: { id: string } }) {
  const router = useRouter();
  const { id } = params;
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  useEffect(() => {
    let active = true;
    const run = async () => {
      try {
        setLoading(true);
        const eq: Equipment = await equipmentApi.getById(id);
        if (!active) return;
        reset({
          equipment_name: eq.equipment_name,
          manufacturer_name: eq.manufacturer_name,
          serial_number: eq.serial_number,
          model_number: eq.model_number,
          category: eq.category,
          customer_name: eq.customer_name,
          installation_location: eq.installation_location,
          installation_date: eq.installation_date?.slice(0,10),
          purchase_date: eq.purchase_date?.slice(0,10),
          warranty_expiry: eq.warranty_expiry?.slice(0,10),
          notes: eq.notes,
        });
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load equipment");
      } finally {
        if (active) setLoading(false);
      }
    };
    run();
    return () => { active = false; };
  }, [id, reset]);

  const onSubmit = async (values: FormValues) => {
    try {
      setSubmitting(true);
      await equipmentApi.update(id, values as any);
      router.push(`/equipment/${id}`);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to update equipment");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <p className="text-gray-600">Loading…</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Card className="max-w-md w-full">
          <CardHeader>
            <CardTitle className="text-red-600">Unable to load</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-gray-700 mb-4">{error}</p>
            <Button className="w-full" onClick={() => router.push("/equipment")}>Back to list</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto">
        <Button variant="ghost" onClick={() => router.push(`/equipment/${id}`)} className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" /> Back
        </Button>

        <Card>
          <CardHeader>
            <CardTitle>Edit Equipment</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm text-gray-700">Equipment Name *</label>
                  <Input {...register("equipment_name")} />
                  {errors.equipment_name && (<p className="text-xs text-red-600 mt-1">{errors.equipment_name.message}</p>)}
                </div>
                <div>
                  <label className="text-sm text-gray-700">Manufacturer *</label>
                  <Input {...register("manufacturer_name")} />
                  {errors.manufacturer_name && (<p className="text-xs text-red-600 mt-1">{errors.manufacturer_name.message}</p>)}
                </div>
                <div>
                  <label className="text-sm text-gray-700">Serial Number *</label>
                  <Input {...register("serial_number")} />
                  {errors.serial_number && (<p className="text-xs text-red-600 mt-1">{errors.serial_number.message}</p>)}
                </div>
                <div>
                  <label className="text-sm text-gray-700">Model Number</label>
                  <Input {...register("model_number")} />
                </div>
                <div>
                  <label className="text-sm text-gray-700">Category</label>
                  <Input {...register("category")} />
                </div>
                <div>
                  <label className="text-sm text-gray-700">Customer Name *</label>
                  <Input {...register("customer_name")} />
                  {errors.customer_name && (<p className="text-xs text-red-600 mt-1">{errors.customer_name.message}</p>)}
                </div>
                <div className="md:col-span-2">
                  <label className="text-sm text-gray-700">Installation Location</label>
                  <Input {...register("installation_location")} />
                </div>
                <div>
                  <label className="text-sm text-gray-700">Installation Date</label>
                  <Input type="date" {...register("installation_date")} />
                </div>
                <div>
                  <label className="text-sm text-gray-700">Purchase Date</label>
                  <Input type="date" {...register("purchase_date")} />
                </div>
                <div>
                  <label className="text-sm text-gray-700">Warranty Expiry</label>
                  <Input type="date" {...register("warranty_expiry")} />
                </div>
                <div className="md:col-span-2">
                  <label className="text-sm text-gray-700">Notes</label>
                  <Input {...register("notes")} />
                </div>
              </div>

              <div className="flex gap-3 justify-end pt-2">
                <Button type="button" variant="outline" onClick={() => router.push(`/equipment/${id}`)}>Cancel</Button>
                <Button type="submit" disabled={submitting}>
                  {submitting ? (<><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving…</>) : "Save"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
