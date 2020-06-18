// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

import { Socket } from "./socket";
import { GoEvent } from "./event";

const reflex = {
    socket: new Socket(window.location.toString()),
    elementID: "",
    async connect() {
        await this.socket.connect()
        this.socket.onmessage = this.onmessage;
    },
    event(event: Event, name: string, args?: any) {
        this.socket.send({ name, event: new GoEvent(event), args });
    },
    onmessage(ev: MessageEvent) {
        // first message is page data
        if (!this.elementID) {
            this.elementID = JSON.parse(ev.data).elementID;
            return;
        }

        const el = document.getElementById(this.elementID);
        if (el) {
            el.innerHTML = ev.data;
        }
    },
};

reflex.connect();

export default reflex;