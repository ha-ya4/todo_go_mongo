function createTodo() {
    const el = document.getElementById("post-todo-form").children
    const todo = {
        id: "",
        title: el[0].value,
        comment: el[1].value
    }
    post("http://localhost:8080", todo).then(res => {
        alert("status:"+res.status)
    }).catch(err => {
        alert(err)
    })
}

function post(url, data) {
    return fetch(url, {
        method: "POST",
        headers:{
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
}