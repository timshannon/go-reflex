// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

import { Socket } from "./socket";


const reflex = {
    socket: new Socket(window.location.toString()),
    async connect() {
        await this.socket.connect()
    },
    event(name: string) {
        console.log("Event: ", name);
    }
};

reflex.connect();

export default reflex;