/**
 * @author: victorystick 
 */

var loader = new THREE.BinaryLoader();
var earlyShipGeom;
var loadModels = function() {
	loader.load('/res/play/models/basic_bin.js', function ( geometry ) {
			geometry.computeBoundingBox();
			geometry.computeBoundingSphere();
			// COMPUTE TANGENTS 
			//geometry.computeTangents();
			earlyShipGeom = geometry;
	})
};

function haxx() {
	var s = local_player.ship;
	s.mesh = newEarlyMesh();
	s.pos = s.mesh.position;
	s.rot = s.mesh.rotation;
	scene.add(s.mesh);
}

var earlyShipMat = new THREE.MeshLambertMaterial( { color: 0x336699, wireframe: false } );

function newEarlyMesh() {
        var mesh = new THREE.Mesh( earlyShipGeom, earlyShipMat );
        mesh.position.x = mesh.position.y = mesh.position.z = 0;
        mesh.rotation.x = mesh.rotation.y = mesh.rotation.z = 0;
        mesh.scale.x = mesh.scale.y = mesh.scale.z = 7;
        mesh.matrixAutoUpdate = false;
        mesh.updateMatrix();
        mesh.matrixAutoUpdate = true;
        return mesh;
};

// SHIP
MP.Ship = function( mesh ) {
	this.mesh = mesh;// || new MP.ErrorBox();
}

MP.Ship.prototype = {
	constructor: MP.Ship,
	update: function() {
		//this.forward
	}
}

// EarlyShip
MP.EarlyShip = function(control) {
/*
	if (MP.Model.EarlyShip) {
		this.mesh = new MP.Model.EarlyShip();
	} else {
		MP.Model.load("EarlyShip", this);
	}
*/
	var geo = new THREE.CylinderGeometry( 0, 7, 30, 3 );
	geo.applyMatrix(new THREE.Matrix4().setRotationFromEuler( new THREE.Vector3( Math.PI / 2, Math.PI, 0 ) ) );
	this.mesh = new THREE.Mesh(geo, new THREE.MeshLambertMaterial( { color: 0x336699, wireframe: false } ));
	this.pos = this.mesh.position;
	this.rot = this.mesh.rotation;
    this.dpos = new THREE.Vector3();
	
	this.id = -1;
	if ( control !== undefined ) {
		this.control = control;
	}
	
	scene.add( this.mesh );
}

MP.EarlyShip.prototype = new MP.Ship();
MP.EarlyShip.prototype = {
	constructor: MP.EarlyShip,
	speed: 0.02,
	mov_speed: 5,
	warp_speed: 30,

	setPos: function( x, y, z ) {
		this.pos.set( x, y, z );
	},

	update : function() {
		if ( this.control !== undefined ) {
			var isDown = this.control.isDown;
			this.move( isDown("SPACE"),
				isDown("CTRL"),
				isDown("SHIFT") );
			this.rotate( isDown("UP"),
				isDown("DOWN"),
				isDown("LEFT"),
				isDown("RIGHT"),
				isDown("ROLL_LEFT"),
				isDown("ROLL_RIGHT"),
				isDown("SHIFT") );			
		}
		this.pos.addSelf(this.dpos);
		
		if(this.target !== undefined) {
			this.mesh.lookAt(this.target);
		}
	},

	move: function( forward, backward, warp ) {
		var speed = (warp)? this.warp_speed: this.mov_speed;
		if ( forward != backward ) {
			if ( backward ) {
				speed /= -2;
			}
			var matrix = new THREE.Matrix4();
			matrix.extractRotation( this.mesh.matrix );
			this.dpos = matrix.multiplyVector3( new THREE.Vector3( 0, 0, 1 ) ).setLength( speed );
			
		} else {
			this.dpos.divideScalar(0);
		}
	},

	rotate: function( up, down, left, right, roll_left, roll_right, warp ) {
		var rotation_matrix = new THREE.Matrix4();
		var speed = this.speed * (warp? 0.05: 1);
		if ( up != down ) {
			if ( down ) {
				rotation_matrix.rotateX(-speed);
			} else {
				rotation_matrix.rotateX(speed);
			}
		}
		if ( left != right ) {
			if ( left ) {
				rotation_matrix.rotateY(speed);
			} else {
				rotation_matrix.rotateY(-speed);
			}
		}
		if ( roll_left != roll_right ) {
			if ( roll_left ) {
				rotation_matrix.rotateZ(-speed);
			} else {
				rotation_matrix.rotateZ(speed);
			}
		}
        this.mesh.matrix.multiplySelf(rotation_matrix);
		this.mesh.rotation.getRotationFromMatrix(this.mesh.matrix);
	},
	toJSON: function() {
		return {"Id": this.id,
			"Pos": this.pos,
			"Rot": this.rot,
			"DPos": this.dpos
		};
	}
};
