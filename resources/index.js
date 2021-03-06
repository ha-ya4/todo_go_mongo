function createTodo() {
    const el = document.getElementById("post-todo-form").children;
    const todo = {
        id: "",
        title: el[0].value,
        comment: el[1].value
    };
    httpPost("http://localhost:8080", todo).then(res => {
        alert("登録成功!");
        location.reload(true);
    }).catch(err => {
        alert(err);
    });
}

function updateTodo(index, todo) {
    const el = document.querySelector(".todo-container"+index).children;
    const t = {
        id: todo.id,
        title: el[0].value,
        comment: el[2].value
    };
    httpPut("http://localhost:8080", t).then(res => {
        alert("変更しました");
        location.reload(true);
    }).catch(err => {
        alert(err);
    });
}

function deleteTodo(todo) {
    const params = new URLSearchParams();
    params.set("id", todo.id);
    const url = "http://localhost:8080/?"+params.toString()
    httpDelete(url).then(res => {
        alert("削除しました")
        location.reload(true);
    }).catch(err => {
        alert(err);
    });
}

function httpPost(url, data) {
    return fetch(url, {
        method: "POST",
        headers:{
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });
}

function httpPut(url, data) {
    return fetch(url, {
        method: "PUT",
        headers:{
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });
}

function httpDelete(url) {
    return fetch(url, {method: "DELETE"});
}