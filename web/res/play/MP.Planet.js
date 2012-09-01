MP.Planet = function(options) {
	this.id = MP.PlanetCount++;
	
	if (options.color !== undefined) {
		this.material = new THREE.MeshLambertMaterial( { color: options.color, wireframe: false } );
	} else {
		this.material = material;
	}
	this.mesh = new THREE.Mesh( new THREE.SphereGeometry(options.size, 10, 16), this.material );
	
	this.mesh.position.x = (options.x)? options.x: 0;
	this.mesh.position.y = (options.y)? options.y: 0;
	this.mesh.position.z = (options.z)? options.z: 0;
	
	scene.add(this.mesh);
}

MP.Planet.prototype = {
	constructor: MP.Planet,

	translate: function(v) {
		this.mesh.translate(v);
	},
	translateX: function(l) {
		this.mesh.translateX(l);
	},
	translateY: function(l) {
		this.mesh.translateY(l);
	},
	translateZ: function(l) {
		this.mesh.translateZ(l);
	}
}

MP.Sun = function(options) {
	var material = new THREE.MeshPhongMaterial( { ambient: 0xffffdd, color: 0xffffff, specular: 0xffffff, shininess: 50, perPixel: true, wireframe: false } );
	this.mesh = new THREE.Mesh( new THREE.SphereGeometry(options.size, 10, 16), material );
	scene.add(this.mesh);
}

MP.Sun.prototype = {
	constructor: MP.Sun
}

MP.PlanetCount = 0;
