export interface Machine {
  id: string;
  name: string;
  status: "online" | "offline";
  IP: string;
  os?: string;
  cpu?: string;
  memory?: string;
  storage?: string;
  uptime?: string;
  users?: string[];
}

export interface Container {
  id: string;
  name: string;
  status: "running" | "stopped" | "paused";
  image: string;
  ports: string;
}

export interface Metrics {
  cpu: {
    current: number;
    average: number;
    peak: number;
  };
  memory: {
    current: number;
    average: number;
    peak: number;
  };
  network: {
    upload: number;
    download: number;
  };
  uptime: number;
}
