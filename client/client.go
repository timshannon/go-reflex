
package client

const Inject = `<script type="text/javascript">
var reflex = (function () {
    'use strict';

    var index = {
        event: function (name) {
            console.log("Event: ", name);
        }
    };

    return index;

}());

</script>`
