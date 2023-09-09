import axios from "axios";
const export_excel = (data) =>{
   return  axios.get('https://localhost:14448/download',data, {
        responseType: 'blob',

    })
        .catch(error => {
            // 处理错误
            console.error(error);
        });
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
    let params = {}
    let res = await export_excel(params);
    // res就是开篇文章那个返回的跟乱码一样的数据
    let blob = res
    let content = [];
    let fileName = 'test.bin'
    //  必须把blob内容放到一个数组里
    content.push(blob);
    downloadFile(content, fileName)
};

export default export_excel