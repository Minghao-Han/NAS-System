const saveToLocalStorage =  (jsonObject) => {
localStorage.setItem("UserInfo",JSON.stringify(jsonObject))
}


export default saveToLocalStorage
