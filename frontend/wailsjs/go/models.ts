export namespace docker {
	
	export class PortMapping {
	    ip: string;
	    privatePort: number;
	    publicPort: number;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new PortMapping(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.privatePort = source["privatePort"];
	        this.publicPort = source["publicPort"];
	        this.type = source["type"];
	    }
	}
	export class ContainerInfo {
	    id: string;
	    shortId: string;
	    names: string[];
	    name: string;
	    image: string;
	    imageId: string;
	    command: string;
	    created: number;
	    state: string;
	    status: string;
	    ports: PortMapping[];
	    sizeRw: number;
	    sizeRootFs: number;
	
	    static createFrom(source: any = {}) {
	        return new ContainerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.shortId = source["shortId"];
	        this.names = source["names"];
	        this.name = source["name"];
	        this.image = source["image"];
	        this.imageId = source["imageId"];
	        this.command = source["command"];
	        this.created = source["created"];
	        this.state = source["state"];
	        this.status = source["status"];
	        this.ports = this.convertValues(source["ports"], PortMapping);
	        this.sizeRw = source["sizeRw"];
	        this.sizeRootFs = source["sizeRootFs"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ContainerStats {
	    id: string;
	    name: string;
	    cpuPercentage: number;
	    memoryUsage: number;
	    memoryLimit: number;
	    memoryPercentage: number;
	    networkRx: number;
	    networkTx: number;
	    blockRead: number;
	    blockWrite: number;
	    pids: number;
	
	    static createFrom(source: any = {}) {
	        return new ContainerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.cpuPercentage = source["cpuPercentage"];
	        this.memoryUsage = source["memoryUsage"];
	        this.memoryLimit = source["memoryLimit"];
	        this.memoryPercentage = source["memoryPercentage"];
	        this.networkRx = source["networkRx"];
	        this.networkTx = source["networkTx"];
	        this.blockRead = source["blockRead"];
	        this.blockWrite = source["blockWrite"];
	        this.pids = source["pids"];
	    }
	}
	
	export class SystemOverview {
	    containers: number;
	    containersRunning: number;
	    containersPaused: number;
	    containersStopped: number;
	    images: number;
	    serverVersion: string;
	    operatingSystem: string;
	    ncpu: number;
	    memTotal: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemOverview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.containers = source["containers"];
	        this.containersRunning = source["containersRunning"];
	        this.containersPaused = source["containersPaused"];
	        this.containersStopped = source["containersStopped"];
	        this.images = source["images"];
	        this.serverVersion = source["serverVersion"];
	        this.operatingSystem = source["operatingSystem"];
	        this.ncpu = source["ncpu"];
	        this.memTotal = source["memTotal"];
	    }
	}
	export class TerminalStartResult {
	    sessionId: string;
	    shell: string;
	
	    static createFrom(source: any = {}) {
	        return new TerminalStartResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.shell = source["shell"];
	    }
	}

}

