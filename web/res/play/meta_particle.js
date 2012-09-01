"use strict";

// GLOBAL VARIABLES
var camera, scene, skybox_scene, renderer, control, composer;
var ships = {}, planets = {};
var local_player;
var ws, ptrSendUpd;
var done = false;
var log, dw;
var updatesRecieved = 0;

// MAIN
jQuery(document).ready(function() {
	log = $("#log");
	dw = $("#debugWindow");
    connect();
    init();
    if ( ws !== undefined ) {
		ptrSendUpd = setInterval(sendUpdate, 1000/10);
    }
    play();
});


function appendLog(msg) {
	if (!fullScreen) {
		var d = log[0]
		var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
		msg.appendTo(log)
		if (doScroll) {
		    d.scrollTop = d.scrollHeight - d.clientHeight;
		}
    }
}


function connect() {
	if (window["WebSocket"]) {
		appendLog($("<div><b>Connecting with WebSockets...</b></div>"))
		var url = window.location.host.split(":");
        ws = new WebSocket("ws://" + url[0] + ":8256/");
        ws.onopen = function(evt) {
		    appendLog($("<div><b>Connection established.</b></div>"))
        }
        ws.onclose = function(evt) {
            appendLog($("<div><b>Connection closed.</b></div>"));
        }
        ws.onmessage = function(evt) {
        	appendLog($("</div>").append(evt.data));
            handleUpdate(evt.data)
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
}

function init() {
	// SCENE
	scene = new THREE.Scene();
	skybox_scene = new THREE.Scene();
	
	// RENDERER
	if (fullScreen) {
		var width = window.innerWidth;
		var height = window.innerHeight;
		var screenRatio=width/height;
   		renderer = new THREE.WebGLRenderer( {antialias: true} );
    	renderer.setSize( window.innerWidth, window.innerHeight );
	} else {
		var screenRatio = 16/9;
		var width = 900;
		var height = width/screenRatio;
   		renderer = new THREE.CanvasRenderer( {antialias: true});
   		renderer.setSize( width, height );
	}

	// CAMERA
    camera = new THREE.PerspectiveCamera( 55, screenRatio, 1, 10000 );
    camera.position.set(2000, 200, 4500);
    scene.add( camera );
    
	/*
	// EFFECTS
	composer = new THREE.EffectComposer( renderer );
	var renderModel = new THREE.RenderPass( scene, camera );
	composer.addPass( renderModel );
	
	if (fullScreen) {
		composer.addPass(new THREE.BloomPass(1.3));
	}
	*/
	
	// lights
	var ambient = new THREE.AmbientLight( 0xffffff );
	ambient.color.setHSV( 0.1, 0.3, 0.1 );
	scene.add( ambient );

/*
	var dirLight = new THREE.DirectionalLight( 0xffffff, 0.125 );
	dirLight.position.set( 0, 0, 1 ).normalize();
	scene.add( dirLight );

	dirLight.color.setHSV( 0.1, 0.725, 0.9 );
*/
	var light = new THREE.PointLight( 0xffffff, 1.5, 7500 );
	light.position.set( 0, 0, 0 );
	scene.add( light );

	light.color.setRGB( 1, 1, 1 );
	
	/*
	// LIGHT
    	scene.add(  new THREE.AmbientLight( 0xffffff ) );
	*/
	
	$("#container").append(renderer.domElement);
	var $canvas = $(renderer.domElement);
	$canvas.attr({"id": "game", "tabindex": "1"});
    
	// LOCAL PLAYER
    control = MP.KeyboardControl.init();
	local_player = new MP.Player( "Player1", camera );
	var ship = new MP.EarlyShip( control );
	local_player.setShip( ship );
	ship.pos.x = 2000;
	ship.pos.y = 200;
	ship.pos.z = 4000;
	ship.mesh.lookAt(new THREE.Vector3(0,0,0));

	// PLANETS
    var planet_options = [
       // { color: 0xffff99, size: 500 },
        { color: 0xffff44, size: 50, x: 5000, y: 0 },
        { color: 0xffff44, size: 50, x: 0, z: 2000 },
        { color: 0xffffaa, size: 5, x: 120, z: 2000 }
    ];
    var planet;
    for ( var i = 0; i < planet_options.length; i++ ) {
        planet = new MP.Planet( planet_options[i] );
        planets[planet.id] = planet;
    }
    
    if (fullScreen) {
	    new MP.Sun( { size:500 } )
	} else {
		planet = new MP.Planet( { color: 0xffff99, size: 500 } );
        planets[planet.id] = planet;
	}

    /* A skybox, for eye-candy and ease of navigation. */
    initSkybox( skybox_scene );
    
	done = true;
}


function initSkybox( scene, urls ) {
    // CREATE A STAR-FILLED SKYBOX

    var texture_placeholder, mesh, materials, geometry;
    var skybox_size = 50000;

    function loadTexture( path ) {
        var texture = new THREE.Texture( texture_placeholder );
        var material = new THREE.MeshBasicMaterial( { map: texture, overdraw: true } );

        var image = new Image();
        image.src = path;
        texture.needsUpdate = true;
        material.map.image = image;

        return material;
    }
    
    //Images made with Starscape
    materials = [
        loadTexture( "/res/play/skybox/red-galaxy-skybox_right1.png" ),
        loadTexture( "/res/play/skybox/red-galaxy-skybox_left2.png" ),
        loadTexture( "/res/play/skybox/red-galaxy-skybox_top3.png" ),
        loadTexture( "/res/play/skybox/red-galaxy-skybox_bottom4.png" ),
        loadTexture( "/res/play/skybox/red-galaxy-skybox_front5.png" ),
        loadTexture( "/res/play/skybox/red-galaxy-skybox_back6.png" )
    ];
    geometry = new THREE.CubeGeometry( skybox_size, skybox_size, skybox_size, 7, 7, 7, materials );
    mesh = new THREE.Mesh( geometry, new THREE.MeshFaceMaterial() );
    mesh.scale.x = -1;
    scene.add(mesh);
}

function logShipinformation() {
    // prints some information about the ship to the log. could use some improvement
    // BUG: The rotations should be from 0 to 360 degrees, not -180 to 180
    var s = "Ship position {", r = "bearing {", d, rad2deg = 180 / Math.PI;
    s += "x: " + parseInt(local_player.ship.pos.x) + ", ";
    s += "y: " + parseInt(local_player.ship.pos.y) + ", ";
    s += "z: " + parseInt(local_player.ship.pos.z) + "} ";
    d = local_player.ship.rot.x * rad2deg; if (d < 0) {}
    r += "x: " + parseInt(d) + "°, ";
    d = local_player.ship.rot.y * rad2deg; if (d < 0) {}
    r += "y: " + parseInt(d) + "°, ";
    d = local_player.ship.rot.z * rad2deg; if (d < 0) {}
    r += "z: " + parseInt(d) + "°}";
    appendLog($("<div>" + s + r + "</div>"));
}

var F1_clicked = false; // a check that keeps the info text from printing every cycle
function play() {
    if (control.isDown("F1")) {
        if (!F1_clicked) {
            F1_clicked = true;
            logShipinformation();
        }
    } else if (F1_clicked) {F1_clicked = false;}
    
	local_player.update();
	
	for (var i in ships) {
		if (ships[i].id != local_player.ship.id) {
			ships[i].update();
		}
	}
	
   	render();
	requestAnimationFrame( play );
}

var enable_skybox = true;
function render() {
	if (done) {
        if (enable_skybox) {
            renderer.render( skybox_scene, camera );
        }
    	renderer.render( scene, camera );
    	//renderer.clear();
    	//composer.render();
    }
}


// SUPPORT FUNCTIONS
function sendUpdate() {
	if (ws.readyState === 3) {
		clearInterval(ptrSendUpd);
	} else {
		if (done) {
			var data = JSON.stringify( local_player );
			if (data) {
				ws.send( data );
			}
		}
	}
}

function handleUpdate( jsonData ) {
	var data = JSON.parse( jsonData );
	var tmp;
	if (!data) {
		return;
	}
	updatesRecieved++;
	//dw.text(jsonData);
	if (data.Assign !== undefined) {
		local_player.ship.id = data.Assign.Id;
		ships[local_player.ship.id] = local_player.ship;
	}
	if (data.Update !== undefined) {
		if (data.Update.Id != local_player.ship.id) {
			var tmp = data.Update;
			if (ships[tmp.Id] === undefined) {
				ships[tmp.Id] = new MP.EarlyShip();
			}
			ships[tmp.Id].pos.copy(tmp.Pos);
			ships[tmp.Id].rot.copy(tmp.Rot);
			ships[tmp.Id].dpos.copy(tmp.DPos);
			//ships[tmp.Id].drot.copy(tmp.DRot);
		}
	}
	if (data.Remove !== undefined) {
		if (ships[data.Remove.Id] !== undefined) {
			scene.remove(ships[data.Remove.Id].mesh);
			delete ships[data.Remove.Id];
		}
	}
}
