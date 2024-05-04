export namespace types {
	
	export class Document {
	    name: string;
	    storageKey: string;
	    originalPath: string;
	    size: number;
	    content: string;
	    metaData: string;
	
	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.storageKey = source["storageKey"];
	        this.originalPath = source["originalPath"];
	        this.size = source["size"];
	        this.content = source["content"];
	        this.metaData = source["metaData"];
	    }
	}

}

