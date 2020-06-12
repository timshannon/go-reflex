// rollup.config.js
import typescript from "@rollup/plugin-typescript";
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
        // terser(),
    ]
};