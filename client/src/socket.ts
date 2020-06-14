// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.


export class Socket {
    public onmessage: ((ev: MessageEvent) => any) | null = null;
    public connection?: WebSocket;
    private manualClose = false;

    constructor(private readonly url: string, public retryPollDuration: number = 5000) { }

    public socketAddress(): string {
        return this.url.replace("http://", "ws://").replace("https://", "wss://");
    }

    public async connect(): Promise<void> {
        let url = this.socketAddress();

        return new Promise((resolve, reject) => {
            this.connection = new WebSocket(url);
            this.connection.onopen = (): void => {
                this.manualClose = false;
                this.connection!.onmessage = (ev: MessageEvent): any => {
                    if (this.onmessage) {
                        this.onmessage(ev);
                    }
                };
                this.connection!.onerror = (event): void => {
                    this.retry();
                };
                // will always retry closed connections until a message is sent from the server to
                // for the client to close the connection themselves.
                this.connection!.onclose = (): void => {
                    if (this.manualClose) {
                        return;
                    }
                    this.retry();
                };
                resolve();
            };

            this.connection.onerror = (event: Event): void => {
                reject(event);
            };

        });
    }

    public async send(data: any): Promise<void> {
        if (!this.connection || this.connection.readyState !== WebSocket.OPEN) {
            await this.connect();
        }

        let msg: string | ArrayBufferLike | Blob | ArrayBufferView;
        if (typeof data === "string" || data instanceof ArrayBuffer || data instanceof Blob) {
            msg = data;
        } else {
            msg = JSON.stringify(data);
        }

        this.connection!.send(msg);
    }

    public close(code?: number, reason?: string): void {
        if (this.connection) {
            this.manualClose = true;
            this.connection.close(code, reason);
        }
    }

    private retry(): void {
        setTimeout(async () => {
            try {
                await this.connect();
            } catch (err) {
                this.retry();
            }
        }, this.retryPollDuration);
    }
}

