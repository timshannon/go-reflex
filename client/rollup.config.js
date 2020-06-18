// rollup.config.js
import typescript from "@rollup/plugin-typescript";
import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
// import { terser } from "rollup-plugin-terser";

export default {
    input: "src/index.ts",
    output: {
        name: "reflex",
        format: "iife",
        file: "client.go",
        banner: `package client\n\nconst Inject = \`<script type="text/javascript">`,
        footer: `</script>\``,
    },
    plugins: [
        typescript(),
        resolve(),
        commonjs(),
        // terser(),  // FIXME: Currently runs after banner and footer are added and breaks
    ]
};