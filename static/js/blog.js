$( document ).ready( function(){

    const like = document.querySelector("a.post-love");
    like.onclick = function(e){
         e.preventDefault()
        fetch(like.getAttribute("href"), {credentials: "same-origin", method: "PUT", redirect: "follow", mode: "same-origin"}).
        then((data) => {
            if (data.status == 307) {
                data.text().then(c => {
                    return window.location.href = c;
                });
            }
            fetch(like.getAttribute("href"), {credentials: "same-origin", method: "GET"}).
            then((res) => {
                res.text().then(c => {
                    like.getElementsByTagName("span")[0].innerText = c;
                })
            });
        });
    }
    document.getElementsByClassName("form-submit")[0].onclick = function(e) {
        e.preventDefault();
        const f =  document.getElementsByClassName("comment-form")[0];
        const formData = new FormData(f);
        if (formData.get("fname").length === 0 || formData.get("email").length === 0 || formData.get("message").length === 0) {
            return;
        }
        fetch(f.action, {
            credentials: "same-origin", 
            method: "POST", 
            redirect: "error",
            headers: {
                'Accept': 'application/json'
              },
              body: new FormData(f)
        })
        .then(() => {
            return fetch("/comments/" + f.querySelector("input[name=blogId]").value, {
                headers: {
                    'Accept': 'application/json'
                  }
            })
        })
        .then(r => {
            return r.text()
        })
        .then(h => {
            f.reset();
            formData.forEach((val, key, fD) => {
                formData.delete(key);
            });
            document.getElementsByClassName("post-comments")[0].innerHTML = h;
        })
    }

});