
const getUserInfo = () => {
    const UserInfo = localStorage.getItem("UserInfo")
    if (UserInfo)
        return JSON.parse(UserInfo)
    else return null
}

export default getUserInfo