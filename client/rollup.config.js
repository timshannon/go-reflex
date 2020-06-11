// rollup.config.js
import typescript from "@rollup/plugin-typescript";
import reflex from "./src/rollup-go-plugin";

export default {
    input: "src/index.ts",
    output: {
        name: "reflex",
        format: "iife",
        file: "client.go",
    },
    plugins: [
        typescript(),
        reflex(),
    ]
};