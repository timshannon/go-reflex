// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

import { Socket } from "./socket";

let socket: Socket | null;

export default {
    async connect() {
        socket = new Socket(window.Location.toString());
        await socket.connect()
    },
    event(name: string) {
        console.log("Event: ", name);
    }
};