import React, { useState } from "react";
import type { Machine } from "@/types/types";
import { MachineCard } from "@/components/machineCard";

export const DashboardPage = () => {
  const [isLoading, setIsLoading] = useState(false);

  // Mock data for machines
  const machines: Machine[] = [
    {
      id: "machine-001",
      name: "Production Server 1",
      status: "online",
      ipAddress: "192.168.1.100"
    },
    {
      id: "machine-002",
      name: "Development Server",
      status: "offline",
      ipAddress: "192.168.1.101"
    },
    {
      id: "machine-003",
      name: "Database Server",
      status: "online",
      ipAddress: "192.168.1.102"
    },
    {
      id: "machine-004",
      name: "Testing Environment",
      status: "offline",
      ipAddress: "192.168.1.103"
    }
  ];

  const handleConnect = (machine: Machine) => {
    console.log(`Connecting to machine: ${machine.name}`);
    setIsLoading(true);
    // Simulate API call
    setTimeout(() => {
      setIsLoading(false);
    }, 2000);
  };

  const handleDisconnect = (machine: Machine) => {
    console.log(`Disconnecting from machine: ${machine.name}`);
    setIsLoading(true);
    // Simulate API call
    setTimeout(() => {
      setIsLoading(false);
    }, 2000);
  };

  const handleEdit = (machine: Machine) => {
    console.log(`Editing machine: ${machine.name}`);
  };

  const handleDelete = (machine: Machine) => {
    console.log(`Deleting machine: ${machine.name}`);
  };

  return (
    <div className="dashboard-container">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-semibold text-foreground mb-2">Machine Dashboard</h1>
          <p className="text-muted-foreground">Manage and monitor your machines</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {machines.map((machine) => (
            <MachineCard
              key={machine.id}
              machine={machine}
              onConnect={handleConnect}
              onDisconnect={handleDisconnect}
              onEdit={handleEdit}
              onDelete={handleDelete}
              isLoading={isLoading}
            />
          ))}
        </div>

        {machines.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No machines found. Add a machine to get started.</p>
          </div>
        )}
      </div>
    </div>
  );
}



