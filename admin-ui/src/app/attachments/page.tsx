'use client'

import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { AttachmentList } from '@/components/attachments/AttachmentList'

export default function AttachmentsPage() {
  return (
    <div className="container mx-auto py-8 space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Attachments & AI Analysis</h1>
          <p className="text-gray-600 mt-2">
            Manage equipment photos and view automated AI analysis results
          </p>
        </div>
      </div>

      <AttachmentList
        showFilters={true}
        showStats={true}
        autoRefresh={true}
        refreshInterval={30000}
      />
    </div>
  )
}