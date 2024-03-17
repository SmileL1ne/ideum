let signInBtnP = document.getElementById("SignInBtnP")
let signUpBtnP = document.getElementById("SignUpBtnP")
let signUpP = document.getElementById("SignUpP")
let signInP = document.getElementById("SignInP")

signUpBtnP.onclick = function(){
    signUpP.style.display = "none"
    signInP.style.display = "flex"

}

signInBtnP.onclick = function(){
    signInP.style.display = "none"
    signUpP.style.display = "flex"

}