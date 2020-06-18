package client

const Inject = `<script type="text/javascript">
var reflex = (function () {
    'use strict';

    /*! *****************************************************************************
    Copyright (c) Microsoft Corporation.

    Permission to use, copy, modify, and/or distribute this software for any
    purpose with or without fee is hereby granted.

    THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
    REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
    AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
    INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
    LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
    OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
    PERFORMANCE OF THIS SOFTWARE.
    ***************************************************************************** */

    function __awaiter(thisArg, _arguments, P, generator) {
        function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
        return new (P || (P = Promise))(function (resolve, reject) {
            function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
            function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
            function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
            step((generator = generator.apply(thisArg, _arguments || [])).next());
        });
    }

    function __generator(thisArg, body) {
        var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
        return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
        function verb(n) { return function (v) { return step([n, v]); }; }
        function step(op) {
            if (f) throw new TypeError("Generator is already executing.");
            while (_) try {
                if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
                if (y = 0, t) op = [op[0] & 2, t.value];
                switch (op[0]) {
                    case 0: case 1: t = op; break;
                    case 4: _.label++; return { value: op[1], done: false };
                    case 5: _.label++; y = op[1]; op = [0]; continue;
                    case 7: op = _.ops.pop(); _.trys.pop(); continue;
                    default:
                        if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                        if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                        if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                        if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                        if (t[2]) _.ops.pop();
                        _.trys.pop(); continue;
                }
                op = body.call(thisArg, _);
            } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
            if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
        }
    }

    // Copyright 2020 Tim Shannon. All rights reserved.
    var Socket = /** @class */ (function () {
        function Socket(url, retryPollDuration) {
            if (retryPollDuration === void 0) { retryPollDuration = 5000; }
            this.url = url;
            this.retryPollDuration = retryPollDuration;
            this.onmessage = null;
            this.manualClose = false;
        }
        Socket.prototype.socketAddress = function () {
            return this.url.replace("http://", "ws://").replace("https://", "wss://");
        };
        Socket.prototype.connect = function () {
            return __awaiter(this, void 0, void 0, function () {
                var url;
                var _this = this;
                return __generator(this, function (_a) {
                    url = this.socketAddress();
                    return [2 /*return*/, new Promise(function (resolve, reject) {
                            _this.connection = new WebSocket(url);
                            _this.connection.onopen = function () {
                                _this.manualClose = false;
                                _this.connection.onmessage = function (ev) {
                                    if (_this.onmessage) {
                                        _this.onmessage(ev);
                                    }
                                };
                                _this.connection.onerror = function (event) {
                                    _this.retry();
                                };
                                // will always retry closed connections until a message is sent from the server to
                                // for the client to close the connection themselves.
                                _this.connection.onclose = function () {
                                    if (_this.manualClose) {
                                        return;
                                    }
                                    _this.retry();
                                };
                                resolve();
                            };
                            _this.connection.onerror = function (event) {
                                reject(event);
                            };
                        })];
                });
            });
        };
        Socket.prototype.send = function (data) {
            return __awaiter(this, void 0, void 0, function () {
                var msg;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!(!this.connection || this.connection.readyState !== WebSocket.OPEN)) return [3 /*break*/, 2];
                            return [4 /*yield*/, this.connect()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            if (typeof data === "string" || data instanceof ArrayBuffer || data instanceof Blob) {
                                msg = data;
                            }
                            else {
                                msg = JSON.stringify(data);
                            }
                            this.connection.send(msg);
                            return [2 /*return*/];
                    }
                });
            });
        };
        Socket.prototype.close = function (code, reason) {
            if (this.connection) {
                this.manualClose = true;
                this.connection.close(code, reason);
            }
        };
        Socket.prototype.retry = function () {
            var _this = this;
            setTimeout(function () { return __awaiter(_this, void 0, void 0, function () {
                var err_1;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            _a.trys.push([0, 2, , 3]);
                            return [4 /*yield*/, this.connect()];
                        case 1:
                            _a.sent();
                            return [3 /*break*/, 3];
                        case 2:
                            err_1 = _a.sent();
                            this.retry();
                            return [3 /*break*/, 3];
                        case 3: return [2 /*return*/];
                    }
                });
            }); }, this.retryPollDuration);
        };
        return Socket;
    }());

    // Copyright 2020 Tim Shannon. All rights reserved.
    // Use of this source code is governed by the MIT license
    // that can be found in the LICENSE file.
    // fields need to match up with event.go
    var GoEvent = /** @class */ (function () {
        function GoEvent(e) {
            this.type = e.type;
            this.altKey = e.altKey;
            this.button = e.button;
            this.buttons = e.buttons;
            this.clientX = e.clientX;
            this.clientY = e.clientY;
            this.ctrlKey = e.ctrlKey;
            this.metaKey = e.metaKey;
            this.movementX = e.movementX;
            this.movementY = e.movementY;
            this.screenX = e.screenX;
            this.screenY = e.screenY;
            this.shiftKey = e.shiftKey;
        }
        return GoEvent;
    }());

    // Copyright 2020 Tim Shannon. All rights reserved.
    var reflex = {
        socket: new Socket(window.location.toString()),
        elementID: "",
        connect: function () {
            return __awaiter(this, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, this.socket.connect()];
                        case 1:
                            _a.sent();
                            this.socket.onmessage = this.onmessage;
                            return [2 /*return*/];
                    }
                });
            });
        },
        event: function (event, name) {
            this.socket.send({ name: name, event: new GoEvent(event) });
        },
        onmessage: function (ev) {
            // first message is page data
            if (!this.elementID) {
                this.elementID = JSON.parse(ev.data).elementID;
                return;
            }
            var el = document.getElementById(this.elementID);
            if (el) {
                el.innerHTML = ev.data;
            }
        },
    };
    reflex.connect();

    return reflex;

}());
</script>`
