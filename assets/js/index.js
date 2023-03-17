let plus = document.querySelector(".plus")
let insertBox = document.querySelector(".insert-box")

let divbody = document.querySelector(".divbody")
let noteBody = document.querySelector(".notebody")

let form = document.querySelector(".insert-box form")

plus.addEventListener("click", (e) => {
    insertBox.classList.toggle("hidden")
})

/* form.addEventListener("submit", (e) => {
    noteBody.innerHTML = divbody.innerHTML
})

console.log(form) */
