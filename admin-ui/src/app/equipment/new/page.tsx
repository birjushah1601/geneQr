"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { equipmentApi } from "@/lib/api/equipment";
import type { RegisterEquipmentRequest } from "@/types";
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
  warranty_months: z.coerce.number().int().min(0).max(120).optional(),
  notes: z.string().optional(),
});

type FormValues = z.infer<typeof schema>;

export default function NewEquipmentPage() {
  const router = useRouter();
  const [submitting, setSubmitting] = useState(false);
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit = async (values: FormValues) => {
    try {
      setSubmitting(true);
      const payload: RegisterEquipmentRequest = {
        ...values,
        created_by: "admin-ui",
      };
      const created = await equipmentApi.register(payload);
      router.push(`/equipment/${created.id}`);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to create equipment");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto">
        <Button variant="ghost" onClick={() => router.push("/equipment")} className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" /> Back
        </Button>

        <Card>
          <CardHeader>
            <CardTitle>Add Equipment</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm text-gray-700">Equipment Name *</label>
                  <Input {...register("equipment_name")} placeholder="Infusion Pump" />
                  {errors.equipment_name && (<p className="text-xs text-red-600 mt-1">{errors.equipment_name.message}</p>)}
                </div>
                <div>
                  <label className="text-sm text-gray-700">Manufacturer *</label>
                  <Input {...register("manufacturer_name")} placeholder="GE Healthcare" />
                  {errors.manufacturer_name && (<p className="text-xs text-red-600 mt-1">{errors.manufacturer_name.message}</p>)}
                </div>
                <div>
                  <label className="text-sm text-gray-700">Serial Number *</label>
                  <Input {...register("serial_number")} placeholder="SN-12345" />
                  {errors.serial_number && (<p className="text-xs text-red-600 mt-1">{errors.serial_number.message}</p>)}
                </div>
                <div>
                  <label className="text-sm text-gray-700">Model Number</label>
                  <Input {...register("model_number")} placeholder="Model X" />
                </div>
                <div>
                  <label className="text-sm text-gray-700">Category</label>
                  <Input {...register("category")} placeholder="Imaging / Monitoring" />
                </div>
                <div>
                  <label className="text-sm text-gray-700">Customer Name *</label>
                  <Input {...register("customer_name")} placeholder="City General Hospital" />
                  {errors.customer_name && (<p className="text-xs text-red-600 mt-1">{errors.customer_name.message}</p>)}
                </div>
                <div className="md:col-span-2">
                  <label className="text-sm text-gray-700">Installation Location</label>
                  <Input {...register("installation_location")} placeholder="Radiology Department, Room 3" />
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
                  <label className="text-sm text-gray-700">Warranty (months)</label>
                  <Input type="number" min={0} max={120} {...register("warranty_months")} />
                </div>
                <div className="md:col-span-2">
                  <label className="text-sm text-gray-700">Notes</label>
                  <Input {...register("notes")} placeholder="Any additional information" />
                </div>
              </div>

              <div className="flex gap-3 justify-end pt-2">
                <Button type="button" variant="outline" onClick={() => router.push("/equipment")}>Cancel</Button>
                <Button type="submit" disabled={submitting}>
                  {submitting ? (<><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving…</>) : "Create"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
