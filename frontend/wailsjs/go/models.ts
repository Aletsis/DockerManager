export namespace container {
	
	export class NetworkInfo {
	    networkName: string;
	    networkId: string;
	    ipAddress: string;
	    gateway: string;
	    macAddress: string;
	    aliases?: string[];
	
	    static createFrom(source: any = {}) {
	        return new NetworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.networkName = source["networkName"];
	        this.networkId = source["networkId"];
	        this.ipAddress = source["ipAddress"];
	        this.gateway = source["gateway"];
	        this.macAddress = source["macAddress"];
	        this.aliases = source["aliases"];
	    }
	}
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
	export class Container {
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
	    networks?: NetworkInfo[];
	    sizeRw: number;
	    sizeRootFs: number;
	    labels?: Record<string, string>;
	    composeProject?: string;
	    composeService?: string;
	    composeWorkingDir?: string;
	    composeConfigFile?: string;
	
	    static createFrom(source: any = {}) {
	        return new Container(source);
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
	        this.networks = this.convertValues(source["networks"], NetworkInfo);
	        this.sizeRw = source["sizeRw"];
	        this.sizeRootFs = source["sizeRootFs"];
	        this.labels = source["labels"];
	        this.composeProject = source["composeProject"];
	        this.composeService = source["composeService"];
	        this.composeWorkingDir = source["composeWorkingDir"];
	        this.composeConfigFile = source["composeConfigFile"];
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
	export class CreateResult {
	    id: string;
	    warnings?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.warnings = source["warnings"];
	    }
	}
	export class CreateSpec {
	    image: string;
	    name?: string;
	    cmd?: string[];
	    ports?: string[];
	    volumes?: string[];
	    env?: string[];
	    restartPolicy?: string;
	    autoStart?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CreateSpec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image = source["image"];
	        this.name = source["name"];
	        this.cmd = source["cmd"];
	        this.ports = source["ports"];
	        this.volumes = source["volumes"];
	        this.env = source["env"];
	        this.restartPolicy = source["restartPolicy"];
	        this.autoStart = source["autoStart"];
	    }
	}
	
	
	export class Stats {
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
	        return new Stats(source);
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

}

export namespace image {
	
	export class DiskUsageSummary {
	    totalImages: number;
	    totalSize: number;
	    danglingCount: number;
	    danglingSize: number;
	    reclaimableSize: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskUsageSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalImages = source["totalImages"];
	        this.totalSize = source["totalSize"];
	        this.danglingCount = source["danglingCount"];
	        this.danglingSize = source["danglingSize"];
	        this.reclaimableSize = source["reclaimableSize"];
	    }
	}
	export class Image {
	    id: string;
	    shortId: string;
	    repository: string;
	    tag: string;
	    repoTags: string[];
	    created: number;
	    size: number;
	    sharedSize: number;
	    containers: number;
	    inUse: boolean;
	    isDangling: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Image(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.shortId = source["shortId"];
	        this.repository = source["repository"];
	        this.tag = source["tag"];
	        this.repoTags = source["repoTags"];
	        this.created = source["created"];
	        this.size = source["size"];
	        this.sharedSize = source["sharedSize"];
	        this.containers = source["containers"];
	        this.inUse = source["inUse"];
	        this.isDangling = source["isDangling"];
	    }
	}
	export class PruneResult {
	    imagesDeleted: string[];
	    spaceReclaimed: number;
	
	    static createFrom(source: any = {}) {
	        return new PruneResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imagesDeleted = source["imagesDeleted"];
	        this.spaceReclaimed = source["spaceReclaimed"];
	    }
	}

}

export namespace network {
	
	export class CreateNetworkSpec {
	    name: string;
	    driver: string;
	    subnet?: string;
	    gateway?: string;
	    ipRange?: string;
	    internal: boolean;
	    attachable: boolean;
	    enableIPv6: boolean;
	    labels?: Record<string, string>;
	    options?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new CreateNetworkSpec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.driver = source["driver"];
	        this.subnet = source["subnet"];
	        this.gateway = source["gateway"];
	        this.ipRange = source["ipRange"];
	        this.internal = source["internal"];
	        this.attachable = source["attachable"];
	        this.enableIPv6 = source["enableIPv6"];
	        this.labels = source["labels"];
	        this.options = source["options"];
	    }
	}
	export class IPAMConfig {
	    driver?: string;
	    subnet?: string;
	    gateway?: string;
	    ipRange?: string;
	
	    static createFrom(source: any = {}) {
	        return new IPAMConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.driver = source["driver"];
	        this.subnet = source["subnet"];
	        this.gateway = source["gateway"];
	        this.ipRange = source["ipRange"];
	    }
	}
	export class NetworkContainerRef {
	    id: string;
	    name: string;
	    state?: string;
	    ipv4Address: string;
	    ipv6Address: string;
	    macAddress: string;
	    endpointId: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkContainerRef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.state = source["state"];
	        this.ipv4Address = source["ipv4Address"];
	        this.ipv6Address = source["ipv6Address"];
	        this.macAddress = source["macAddress"];
	        this.endpointId = source["endpointId"];
	    }
	}
	export class Network {
	    id: string;
	    shortId: string;
	    name: string;
	    driver: string;
	    scope: string;
	    internal: boolean;
	    attachable: boolean;
	    enableIPv6: boolean;
	    ipam: IPAMConfig[];
	    containers: NetworkContainerRef[];
	    containersCount: number;
	    labels: Record<string, string>;
	    options: Record<string, string>;
	    created: string;
	    isDefault: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Network(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.shortId = source["shortId"];
	        this.name = source["name"];
	        this.driver = source["driver"];
	        this.scope = source["scope"];
	        this.internal = source["internal"];
	        this.attachable = source["attachable"];
	        this.enableIPv6 = source["enableIPv6"];
	        this.ipam = this.convertValues(source["ipam"], IPAMConfig);
	        this.containers = this.convertValues(source["containers"], NetworkContainerRef);
	        this.containersCount = source["containersCount"];
	        this.labels = source["labels"];
	        this.options = source["options"];
	        this.created = source["created"];
	        this.isDefault = source["isDefault"];
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
	
	export class PruneResult {
	    networksDeleted: string[];
	
	    static createFrom(source: any = {}) {
	        return new PruneResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.networksDeleted = source["networksDeleted"];
	    }
	}

}

export namespace system {
	
	export class Overview {
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
	        return new Overview(source);
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

}

export namespace terminal {
	
	export class StartResult {
	    sessionId: string;
	    rows: number;
	    cols: number;
	
	    static createFrom(source: any = {}) {
	        return new StartResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.rows = source["rows"];
	        this.cols = source["cols"];
	    }
	}

}

export namespace volume {
	
	export class ContainerRef {
	    id: string;
	    name: string;
	    state: string;
	    destination: string;
	    rw: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ContainerRef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.state = source["state"];
	        this.destination = source["destination"];
	        this.rw = source["rw"];
	    }
	}
	export class DiskUsageSummary {
	    totalVolumes: number;
	    totalSize: number;
	    danglingCount: number;
	    danglingSize: number;
	    reclaimableSize: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskUsageSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalVolumes = source["totalVolumes"];
	        this.totalSize = source["totalSize"];
	        this.danglingCount = source["danglingCount"];
	        this.danglingSize = source["danglingSize"];
	        this.reclaimableSize = source["reclaimableSize"];
	    }
	}
	export class PruneResult {
	    volumesDeleted: string[];
	    spaceReclaimed: number;
	
	    static createFrom(source: any = {}) {
	        return new PruneResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.volumesDeleted = source["volumesDeleted"];
	        this.spaceReclaimed = source["spaceReclaimed"];
	    }
	}
	export class Volume {
	    name: string;
	    driver: string;
	    mountpoint: string;
	    createdAt: string;
	    labels: Record<string, string>;
	    scope: string;
	    size: number;
	    inUse: boolean;
	    containers: ContainerRef[];
	
	    static createFrom(source: any = {}) {
	        return new Volume(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.driver = source["driver"];
	        this.mountpoint = source["mountpoint"];
	        this.createdAt = source["createdAt"];
	        this.labels = source["labels"];
	        this.scope = source["scope"];
	        this.size = source["size"];
	        this.inUse = source["inUse"];
	        this.containers = this.convertValues(source["containers"], ContainerRef);
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

}

