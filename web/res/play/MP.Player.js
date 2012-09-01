MP.Player = function(name, camera, ship) {
	this.name = name || "NewPlayer";
	this.camera = camera;
	this.ship = ship;
    this._matrix4 = new THREE.Matrix4();
}

MP.Player.prototype = {
	constructor: MP.Player,
	update: function() {
			if ( this.ship !== undefined ) {
				this.ship.update();
				this.dragCamera();
                //this.followCamera();
			}
		},
	setShip: function( ship ) {
			this.ship = ship;
		},
	toJSON: function() {
		if ( this.ship.control.changed ) {
			this.ship.control.changed = false;
			return { "Update": this.ship.toJSON() };
		}
		return null;
	},
	dragCamera: function() {
		var ny = new THREE.Vector3();
		ny.sub(this.camera.position, this.ship.pos)
		if (ny.length() > 300) {
			ny.setLength(300);
		}
        
		this.camera.position.add(this.ship.pos, ny);
		this.camera.lookAt(this.ship.pos);
	},
    followCamera: function() {
        // Trying to add a follow-cam.
        var deg2rad = Math.PI / 180;
        var distance = 100; 
        var cameraPosition = new THREE.Vector3( 0, 0.1, -1 ); // camera vector
        
        var matrix = this._matrix4.extractRotation( this.ship.mesh.matrix );
        cameraPosition = matrix.multiplyVector3( cameraPosition )
        cameraPosition.setLength( distance );
		this.camera.position.add(this.ship.pos, cameraPosition);
        
        // this.camera.matrix.extractRotation( this.ship.mesh.matrix );
        // cameraPosition = this.camera.matrix.multiplyVector3( cameraPosition )
        // cameraPosition.setLength( distance );
		// this.camera.position.add(this.ship.pos, cameraPosition);
        
        this.camera.rotation.z = -this.ship.rot.z;
        this.camera.rotation.y = this.ship.rot.y + 180*deg2rad;
        this.camera.rotation.x = this.ship.rot.x;
        // this.camera.updateMatrix(); // seems to be called at a later time, unneccessary here
    }
}
