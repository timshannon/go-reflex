// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// fields need to match up with event.go
export class GoEvent {
    public type: string;
    public altKey: boolean;
    public button: number;
    public buttons: number;
    public clientX: number;
    public clientY: number;
    public ctrlKey: boolean;
    public metaKey: boolean;
    public movementX: number;
    public movementY: number;
    public screenX: number;
    public screenY: number;
    public shiftKey: boolean;
    constructor(e: any) {
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
}
