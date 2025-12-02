import React, { useEffect, useState } from "react";
import type { Machine } from "@/types/types";
import { MachineCard } from "@/components/machineCard";
import axios from "axios";

export const DashboardPage = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [machines, setMachines] = useState<Machine[]>([]);
  // Mock data for machines
  useEffect(() => {
    async function FetchMachines() {
      try {
        const res = await axios.get("http://localhost:8080/api/user/machine/usable", {
          headers: {
            Authorization: localStorage.getItem("token"),
          },
        });
        console.log(res.data);
        const machinesData: Machine[] = [];
        res.data.usable_machines.forEach((machine: any) => {
          machinesData.push({
            id: machine.ID,
            name: machine.Name,
            status: "online",
            IP: machine.IP,
          });
        });
        setMachines(machinesData);
      } catch (error) {
        console.error("There was an error fetching the machines!", error);
      }
    }
    FetchMachines();
  }, []);

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
};
