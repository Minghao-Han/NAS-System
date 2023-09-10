import axios from "axios";
const getPort = (data) =>{
    let dsport
    console.log("get")
    axios.get('https://localhost:443/user/file/large', {
        headers: {
            'Content-Type': 'application/json',
            'token':localStorage.getItem("token")
        }
    })
        .then(response => {
            // 处理响应
          dsport = response.data.dsPort
            console.log("dsport")
            console.log(dsport)
            console.log(response)
            getData(dsport,data)

        })
        .catch(error => {
            // 处理错误
            console.error(error);
        });

    return dsport

}

const getData = (dsport,data) => {

    console.log("port")
    console.log(dsport)
     axios.get('https://localhost:'+dsport+'/download', {


        headers: {
            'token':localStorage.getItem("token"),
            'filePath':'\\test.jpg',
        },
        responseType:'blob'
    })
        .then(response => {
            // 处理响应
            console.log("j结果")
            console.log(response)
            let blob = response
            let content = [];
            let fileName = 'test.bin'
            //  必须把blob内容放到一个数组里
            content.push(blob);
            downloadFile(content, fileName)

        })
        .catch(error => {
            // 处理错误
            console.error(error);
        });
     return null
}

const downloadFile = (data, fileName) => {
    // new Blob 实例化文件流
    const url = window.URL.createObjectURL(new Blob(data))
    const link = document.createElement('a')
    link.style.display = "none"
    link.href = url
    link.setAttribute('download', fileName)
    document.body.appendChild(link)
    link.click()
    //下载完成移除元素
    document.body.removeChild(link)
    //释放掉blob对象
    window.URL.revokeObjectURL(url)
};

const exportData = async () => {
    let data="/test.jpg"
    let res = await getPort(data);
    // res就是开篇文章那个返回的跟乱码一样的数据

};

export default exportData