MP.KeyboardControl = {
	key : {
		_map: {
			65: "LEFT",
			68: "RIGHT",
			81: "ROLL_LEFT",
			69: "ROLL_RIGHT",
			83: "DOWN",
			87: "UP",
			32: "SPACE",
			17: "CTRL",
			16: "SHIFT", 
            112: "F1"
		},
		_pressed : {
			"LEFT": false,
			"RIGHT": false,
			"ROLL_LEFT": false,
			"ROLL_RIGHT": false,
			"DOWN": false,
			"UP": false,
			"SPACE": false,
			"CTRL": false,
			"SHIFT": false,
            "F1": false
		}
	},
	isDown : function(key) {
		return MP.KeyboardControl.key._pressed[key];
	},
	_notIn : function(keyCode) {
		var ret = false;
		jQuery.each(MP.KeyboardControl.key._map, function(key, value) {
			if(key == keyCode) {
				ret = true;
				return false;
			}
		});
		return ret;
	},
	_game : [],
	handlers : {
		keydown: function(event) {
			var ctrl = MP.KeyboardControl;
			if (ctrl._notIn(event.which)) {
				if (ctrl.key._pressed[ctrl.key._map[event.which]] == false) {
					ctrl.key._pressed[ctrl.key._map[event.which]] = true;
					ctrl.changed = true;
				}
			}
			return false;
		},
		keyup: function(event) {
			var ctrl = MP.KeyboardControl;
			if (ctrl._notIn(event.which)) {
				if (ctrl.key._pressed[ctrl.key._map[event.which]] == true) {
					ctrl.key._pressed[ctrl.key._map[event.which]] = false;
					ctrl.changed = true;
				}
			}
			return false;
		},
		mousedown: function(event) {
			var ctrl = MP.KeyboardControl;
			ctrl._game.on("mousemove", ctrl.handlers._mouse);
		},
		mouseup: function(event) {
			var ctrl = MP.KeyboardControl;
			ctrl._game.off("mousemove", ctrl.handlers._mouse);
		},
		_mouse : function(event) {
		}
	},
	changed : true,
	init : function() {
		var ctrl = MP.KeyboardControl;
		ctrl._game = $("#game");
		ctrl._game.on(ctrl.handlers);
		ctrl._game.off("_mouse");
		ctrl._game.focus();
		return this;
	}
};
