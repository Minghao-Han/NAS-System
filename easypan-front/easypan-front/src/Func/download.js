import axios from "axios";

const getPort = (data) =>{
    let dsport
    console.log("get")
    axios.get('https://localhost:443/user/file/large', {
        headers: {
            'Content-Type': 'application/json',
            'token':localStorage.getItem("token"),
            'path':'\\',
            'filename':'test.mp4',
            'size':'34'
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
            'filePath':'/test.mp4',
            'size':'34'
        },
        responseType:'blob'
    })
        .then(response => {
            // 处理响应
            console.log("j结果")
            console.log(response)
            let content = [];
            let fileName = 'test.mp4'


            // //读取文件
            // // 1.创建 FileReader 对象
            // const fileReader = new FileReader()
            // // 2.开始读取指定的Blob中的内容。一旦完成，result属性中将包含一个字符串以表示所读取的文件内容。
            // fileReader.readAsText(response.data) // 3.以字符串的形式读取出来
            // fileReader.onload = () => {
            //     // 4.该事件在读取操作完成时触发。注意：此时this指向fileReader
            //
            //     let result = JSON.parse(result) //获取的结果根据后端返回的数据类型选用json.parse
            //     if (result.code !== 0) {
            //         //如果读取失败进行响应的操作或提示
            //         return this.$message.error('文件读取失败')
            //     }
            // }
//导出文件

            content.push(response.data);
            console.log("content.length")
            console.log(content.length)
            const url = window.URL.createObjectURL(new Blob(content))
            const link = document.createElement('a')
            link.style.display = "none"
            link.href = url
            link.setAttribute('download', fileName)
            document.body.appendChild(link)
            link.click()
            setTimeout(function() {
                window.URL.revokeObjectURL(url);
                document.body.removeChild(link);
            }, 100);
            //content.push(blob);
            //downloadFile(content, fileName)

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
    setTimeout(function() {
        window.URL.revokeObjectURL(url);
        document.body.removeChild(link);
    }, 100);
};

const exportData = async () => {
    let data="/test.mp4"
    let res = await getPort(data);
    // res就是开篇文章那个返回的跟乱码一样的数据

};

export default exportData