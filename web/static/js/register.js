let signInBtn = document.getElementById("SignInBtn")
let signUpBtn = document.getElementById("SignUpBtn")
let signUp = document.getElementById("SignUp")
let signIn = document.getElementById("SignIn")

signUpBtn.onclick = function(){
    signUp.style.display = "none"
    signIn.style.display = "flex"

}

signInBtn.onclick = function(){
    signIn.style.display = "none"
    signUp.style.display = "flex"

}