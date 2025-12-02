import React from "react";
import type { Machine } from "@/types/types";

interface MachineCardProps {
  machine: Machine;
  onConnect?: (machine: Machine) => void;
  onDisconnect?: (machine: Machine) => void;
  onEdit?: (machine: Machine) => void;
  onDelete?: (machine: Machine) => void;
  isLoading?: boolean;
  className?: string;
}

export function MachineCard({
  machine,
  onConnect,
  onDisconnect,
  onEdit,
  onDelete,
  isLoading = false,
  className = ""
}: MachineCardProps) {
  const isOnline = machine.status === "online";

  return (
    <div className={`bg-card border border-border rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow ${className}`}>
      <div className="space-y-4">
        {/* Header with machine name and status */}
        <div className="flex items-center justify-between">
          <h3 className="text-xl font-semibold text-card-foreground">{machine.name}</h3>
          <div className={`px-3 py-1 rounded-full text-sm font-medium ${
            isOnline
              ? "bg-green-100 text-green-700 border border-green-200"
              : "bg-muted text-muted-foreground border border-border"
          }`}>
            {machine.status}
          </div>
        </div>

        {/* Machine details */}
        <div className="space-y-2">
          <div className="text-sm text-muted-foreground">
            <span className="font-medium text-card-foreground">Machine ID:</span> {machine.id}
          </div>
          <div className="text-sm text-muted-foreground">
            <span className="font-medium text-card-foreground">IP Address:</span> {machine.ipAddress}
          </div>
        </div>

        {/* Action buttons */}
        <div className="flex gap-2 pt-2">
          {isOnline && (
            <button
              onClick={() => onDisconnect?.(machine)}
              disabled={isLoading}
              className="px-4 py-2 bg-destructive hover:bg-destructive/90 text-destructive-foreground text-sm rounded border border-destructive disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
            >
              {isLoading ? "Disconnecting..." : "Disconnect"}
            </button>
          )}

          <button
            onClick={() => onEdit?.(machine)}
            disabled={isLoading}
            className="px-4 py-2 bg-secondary hover:bg-secondary/80 text-secondary-foreground text-sm rounded border border-border disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
          >
            Edit
          </button>

          <button
            onClick={() => onDelete?.(machine)}
            disabled={isLoading}
            className="px-4 py-2 bg-card hover:bg-accent text-destructive text-sm rounded border border-border disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}
