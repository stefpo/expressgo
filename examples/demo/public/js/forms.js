	$e=function(id) {
		return document.getElementById(id);
	}
	
	function GfForm( e ) {
		ret = new Object()
		ret.form = e;
		ret.isValue = false;
		ret.cssRequiredField = "required-field";
		ret.cssMissingField = "missing-field";
		
		ret.toObject = function () {
			var i;
			var fo = new Object();
			var form = this.form;
			
			this.isValid = 1;
			for ( i = 0; i < form.elements.length ; i++) {
				e = form.elements[i];
				if (e.getAttribute("gf-required") != undefined ) { addClass (e, this.cssRequiredField) }
				removeClass (e, this.cssMissingField);
				id = e.getAttribute("name");
				if (id !=null && id != "") {
					if (e.tagName == "INPUT" ) {
						
						if (e.type == "checkbox") {
							fo[id] = e.checked ? 1 : 0;
						}
						else if (e.type == "radio") {
							if (e.checked) 	fo[id] = e.value;
						}
						
						else {
							if (e.getAttribute("gf-required") != undefined && e.value =="") {
								this.isValid = 0;
								addClass (e, this.cssMissingField); 
							}
							if (e.getAttribute("gf-type") == "float" ) {
								e.value = e.value.trim();
								fo[id] = parseFloat(e.value);
								if (isNaN (fo[id] )) {
									addClass (e, this.cssMissingField); 
									this.isValid = 0;
								}
							} 
							else if (e.getAttribute("gf-type") == "date" ) {
								e.value = e.value.trim();
								fo[id] = new Date(e.value);
								
								if (isNaN (fo[id].getTime() )) {
									addClass (e, this.cssMissingField); 
									this.isValid = 0;
								}
							} 						
							else {
								fo[id] = e.value;
							}
						}
					}
					else if (e.tagName == "SELECT" ) {
						fo[id] = e.value;
					}
				}
			}
			fo['$valid']=this.isValid;
			return (fo);
		}
		
		ret.fill = function ( fo ) {
			var form=this.form;
			var i;
			for ( i = 0; i < form.elements.length ; i++) {
				e = form.elements[i];
				id = e.getAttribute("name");
				if (id !=null && id != "") {
					if (fo[id] != undefined) {
						
						if (e.tagName == "INPUT" ) {
							
							if (e.type == "checkbox") {
								e.checked = fo[id] == 0 ? false : true;
							}
							else if (e.type == "radio") {
								if (fo[id] == e.value) e.checked = true;
							}
							
							else {
								e.value = fo[id];
							}
						}
						else if (e.tagName == "SELECT" ) {
							e.value = fo[id];
						}
					}
				}
			}
			return (fo);
		}	
		x = ret.toObject();
		
		return ret;
	}

	function addClass(e, classname) {
		var css = e.className.split(" ");
		var i;
		var newcss ="";
		for (i = 0; i< css.length; i++) {
			if ( css[i] != "" && css[i] != classname && css[i] != undefined) {
				if ( newcss != "") newcss = newcss + " ";
				newcss = newcss + css[i];
			}
		}
		if ( newcss != "") newcss = newcss + " ";
		newcss = newcss + classname;
		e.className = newcss;
	}
	
	function removeClass(e, classname) {
		var css = e.className.split(" ");
		var i;
		var newcss = "";
		for (i = 0; i< css.length; i++) {
			if ( css[i] != "" && css[i] !=classname && css[i] !== undefined) {
				if ( newcss != "") newcss = newcss + " ";
				newcss = newcss + css[i];
			}
		}
		e.className = newcss;
	}
	

	function getElementAttributes(e) {
		var i;
		var r = Object() ; 
		for ( i = 0; i < e.attributes.length; i++ ) {
			r[e.attributes[i].name] = e.attributes[i].value;
		}
		return r;
	}






	
