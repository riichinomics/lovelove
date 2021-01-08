import { Api } from "./rpc/Api";

const api = new Api({url: "localhost:6482"});
api.init().then(() => {
    api.lovelove.sayHello({
        name: "test"
    }).then((response) => {
        console.log("test", response)
    });
});
