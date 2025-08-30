export namespace main {
	
	export class FinalCdResult {
	    InputFileName: string;
	    OutputFilePath: string;
	    Success: boolean;
	    Stdout: string;
	    Err: string;
	
	    static createFrom(source: any = {}) {
	        return new FinalCdResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.InputFileName = source["InputFileName"];
	        this.OutputFilePath = source["OutputFilePath"];
	        this.Success = source["Success"];
	        this.Stdout = source["Stdout"];
	        this.Err = source["Err"];
	    }
	}

}

