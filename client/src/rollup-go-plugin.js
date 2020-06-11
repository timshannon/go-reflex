export default function goreflex(options = {}) {
    return {
        name: "go-reflex-client",
        banner() {
            return `package client\n\nconst Inject = \`<script type="text/javascript">`;
        },
        footer() {
            return `</script>\``
        },
    };
}