import React from "react";
import { useNavigate } from "react-router-dom";
import { Server, Box, BarChart3, Settings, Trash2, Power, PowerOff } from "lucide-react";
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
  const navigate = useNavigate();
  const isOnline = machine.status === "online";

  const handleViewMachine = () => {
    navigate(`/machine/${machine.id}`);
  };

  const handleViewContainers = () => {
    navigate(`/machine/${machine.id}?view=containers`);
  };

  const handleViewMetrics = () => {
    navigate(`/machine/${machine.id}?view=metrics`);
  };

  return (
    <div className={`bg-card border border-border rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow ${className}`}>
      <div className="space-y-4">
        {/* Header with machine name and status */}
        <div className="flex items-center justify-between">
          <h3 className="text-xl font-semibold text-card-foreground cursor-pointer hover:text-blue-600" onClick={handleViewMachine}>
            {machine.name}
          </h3>
          <div className={`px-3 py-1 rounded-full text-sm font-medium ${
            isOnline
              ? "bg-green-100 text-green-700 border border-green-200"
              : "bg-red-100 text-red-700 border border-red-200"
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
            <span className="font-medium text-card-foreground">IP Address:</span> {machine.IP}
          </div>
        </div>

        {/* Quick access buttons */}
        <div className="flex gap-2 pt-2">
          <button
            onClick={handleViewMachine}
            className="flex items-center gap-1 px-3 py-1.5 bg-blue-100 hover:bg-blue-200 text-blue-700 text-xs rounded border border-blue-200 transition-colors font-medium"
            title="View Machine Details"
          >
            <Server className="h-3 w-3" />
            Details
          </button>

          <button
            onClick={handleViewContainers}
            className="flex items-center gap-1 px-3 py-1.5 bg-purple-100 hover:bg-purple-200 text-purple-700 text-xs rounded border border-purple-200 transition-colors font-medium"
            title="View Containers"
          >
            <Box className="h-3 w-3" />
            Containers
          </button>

          <button
            onClick={handleViewMetrics}
            className="flex items-center gap-1 px-3 py-1.5 bg-green-100 hover:bg-green-200 text-green-700 text-xs rounded border border-green-200 transition-colors font-medium"
            title="View Metrics"
          >
            <BarChart3 className="h-3 w-3" />
            Metrics
          </button>
        </div>

        {/* Action buttons */}
        <div className="flex gap-2 pt-2 border-t border-border">
          {isOnline ? (
            <button
              onClick={() => onDisconnect?.(machine)}
              disabled={isLoading}
              className="flex items-center gap-1 px-3 py-2 bg-destructive hover:bg-destructive/90 text-destructive-foreground text-sm rounded border border-destructive disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
            >
              <PowerOff className="h-3 w-3" />
              {isLoading ? "Disconnecting..." : "Disconnect"}
            </button>
          ) : (
            <button
              onClick={() => onConnect?.(machine)}
              disabled={isLoading}
              className="flex items-center gap-1 px-3 py-2 bg-green-600 hover:bg-green-700 text-white text-sm rounded border border-green-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
            >
              <Power className="h-3 w-3" />
              {isLoading ? "Connecting..." : "Connect"}
            </button>
          )}

          <button
            onClick={() => onEdit?.(machine)}
            disabled={isLoading}
            className="flex items-center gap-1 px-3 py-2 bg-secondary hover:bg-secondary/80 text-secondary-foreground text-sm rounded border border-border disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
          >
            <Settings className="h-3 w-3" />
            Edit
          </button>

          <button
            onClick={() => onDelete?.(machine)}
            disabled={isLoading}
            className="flex items-center gap-1 px-3 py-2 bg-card hover:bg-accent text-destructive text-sm rounded border border-border disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
          >
            <Trash2 className="h-3 w-3" />
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}
