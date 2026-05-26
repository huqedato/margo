export namespace main {
	
	export class FileDetails {
	    filename: string;
	    path: string;
	    size: number;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new FileDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.content = source["content"];
	    }
	}
	export class ReturnResult {
	    filename: string;
	    path: string;
	    html: string;
	
	    static createFrom(source: any = {}) {
	        return new ReturnResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.path = source["path"];
	        this.html = source["html"];
	    }
	}

}

