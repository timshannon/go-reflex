// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

import { Socket } from "./socket";
import { GoEvent } from "./event";


const reflex = {
    socket: new Socket(window.location.toString()),
    async connect() {
        await this.socket.connect()
        this.socket.onmessage = this.onmessage;
    },
    event(event: Event, name: string) {
        this.socket.send({ name, event: new GoEvent(event) });
    },
    onmessage(ev: MessageEvent) {
        console.log(ev);
    },
};

reflex.connect();

export default reflex;