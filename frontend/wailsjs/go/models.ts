export namespace main {
	
	export class AppInfo {
	    // Go type: struct { ProductName string "json:\"productName\""; ProductVersion string "json:\"productVersion\""; CompanyName string "json:\"companyName\""; Copyright string "json:\"copyright\""; Description string "json:\"description\""; License string "json:\"license\"" }
	    info: any;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.info = this.convertValues(source["info"], Object);
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

