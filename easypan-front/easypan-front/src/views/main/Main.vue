<template>
  <div>
    <div class="top">
      <div class="top-op">
        <div class="btn">
          <el-dialog :title="upload.title" v-model="upload.open" width="400px" append-to-body>
            <el-upload
                ref="uploadRef"
                :limit="1"
                :headers="upload.headers"
                :action="upload.url"
                :disabled="upload.isUploading"
                :on-progress="handleFileUploadProgress"
                :on-success="handleFileSuccess"
                :auto-upload="false"
                drag>
              <el-icon class="el-icon--upload">
                <upload-filled/>
              </el-icon>
              <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
              <!--              <template #tip>-->
              <!--                <div class="el-upload__tip text-center">-->
              <!--                  &lt;!&ndash;                  <div class="el-upload__tip">&ndash;&gt;-->
              <!--                  &lt;!&ndash;                    <el-checkbox v-model="upload.updateSupport"/>&ndash;&gt;-->
              <!--                  &lt;!&ndash;                    是否更新已经存在的用户数据&ndash;&gt;-->
              <!--                  &lt;!&ndash;                  </div>&ndash;&gt;-->
              <!--                  &lt;!&ndash;                  <span>仅允许导入xls、xlsx格式文件。</span>&ndash;&gt;-->
              <!--                  <el-link type="primary" :underline="false" style="font-size: 12px" @click="importTemplate">下载模板-->
              <!--                  </el-link>-->
              <!--                </div>-->
              <!--              </template>-->
            </el-upload>
            <template #footer>
              <div class="dialog-footer">
                <el-button type="primary" @click="submitFileForm">确 定</el-button>
                <el-button @click="upload.open = false">取 消</el-button>
              </div>
            </template>
          </el-dialog>

          <!--          <el-upload-->
          <!--              :show-file-list="false"-->
          <!--              :with-credentials="true"-->
          <!--              :multiple="true"-->
          <!--              :http-request="addFile"-->
          <!--              :beforeUpload="beforeUpload"-->
          <!--              :accept="fileAccept"-->
          <!--          >-->
          <el-button type="primary" @click="upload.open = true">
            <span class="iconfont icon-upload"></span>
            上传
          </el-button>

          <!--          </el-upload>-->
        </div>
        <el-button
            type="primary"
            @click="downloadFile"
        ><span class="iconfont icon-download">下载</span>
        </el-button>
        <!--        <el-button type="primary" @click="downloadFile">-->
        <!--          <span class="iconfont icon-download"></span>-->
        <!--          下载-->
        <!--        </el-button>-->
        <el-button type="success" @click="newFolder" v-if="category == 'all'">
          <span class="iconfont icon-folder-add"></span>
          新建文件夹
        </el-button>
        <el-button
            @click="delFileBatch"
            type="danger"
            :disabled="selectFileIdList.length == 0"
        >
          <span class="iconfont icon-del"></span>
          批量删除
        </el-button>
        <el-button
            @click="moveFolderBatch"
            type="warning"
            :disabled="selectFileIdList.length == 0"
        >
          <span class="iconfont icon-move"></span>
          批量移动
        </el-button>
        <div class="search-panel">
          <el-input
              clearable
              placeholder="输入文件名搜索"
              v-model="fileNameFuzzy"
              @keyup.enter="search"
          >
            <template #suffix>
              <i class="iconfont icon-search" @click="search"></i>
            </template>
          </el-input>
        </div>
        <div class="iconfont icon-refresh" @click="loadDataList"></div>
      </div>
      <!--导航-->
      <Navigation ref="navigationRef" @navChange="navChange"></Navigation>
    </div>
    <div class="file-list" v-if="tableData.list && tableData.list.length > 0">
      <Table
          ref="dataTableRef"
          :columns="columns"
          :showPagination="true"
          :dataSource="tableData"
          :fetch="loadDataList"
          :initFetch="false"
          :options="tableOptions"
          @rowSelected="rowSelected"
      >
        <template #fileName="{ index, row }">
          <div
              class="file-item"
              @mouseenter="showOp(row)"
              @mouseleave="cancelShowOp(row)"
          >
            <template
                v-if="(row.fileType == 3 || row.fileType == 1) && row.status == 2"
            >
              <icon :cover="row.fileCover" :width="32"></icon>
            </template>
            <template v-else>
              <icon v-if="row.folderType == 0" :fileType="row.fileType"></icon>
              <icon v-if="row.folderType == 1" :fileType="0"></icon>
            </template>
            <span class="file-name" v-if="!row.showEdit" :title="row.fileName">
              <span @click="preview(row)">{{ row.fileName }}</span>
              <span v-if="row.status == 0" class="transfer-status">转码中</span>
              <span v-if="row.status == 1" class="transfer-status transfer-fail"
              >转码失败</span
              >
            </span>
            <div class="edit-panel" v-if="row.showEdit">
              <el-input
                  v-model.trim="row.fileNameReal"
                  ref="editNameRef"
                  :maxLength="190"
                  @keyup.enter="saveNameEdit(index)"
              >
                <template #suffix>{{ row.fileSuffix }}</template>
              </el-input>
              <span
                  :class="[
                  'iconfont icon-right1',
                  row.fileNameReal ? '' : 'not-allow',
                ]"
                  @click="saveNameEdit(index)"
              ></span>
              <span
                  class="iconfont icon-error"
                  @click="cancelNameEdit(index)"
              ></span>
            </div>
            <span class="op">
              <template v-if="row.showOp && row.fileId && row.status == 2">
                <span class="iconfont icon-share1" @click="share(row)"
                >分享</span
                >
                <!--                <span-->
                <!--                    class="iconfont icon-download"-->
                <!--                    @click="download(row)"-->
                <!--                    v-if="row.folderType == 0"-->
                <!--                >下载</span-->
                <!--                >-->
                <span class="iconfont icon-del" @click="delFile(row)"
                >删除</span
                >
                <span
                    class="iconfont icon-edit"
                    @click.stop="editFileName(index)"
                >重命名</span
                >
                <span class="iconfont icon-move" @click="moveFolder(row)"
                >移动</span
                >
              </template>
            </span>
          </div>
        </template>
        <template #fileSize="{ index, row }">
          <span v-if="row.fileSize">
            {{ proxy.Utils.size2Str(row.fileSize) }}</span
          >
        </template>
      </Table>
    </div>
    <div class="no-data" v-else>
      <div class="no-data-inner">
        <Icon iconName="no_data" :width="120" fit="fill"></Icon>
        <div class="tips">当前目录为空，上传你的第一个文件吧</div>
        <div class="op-list">
          <el-upload
              :show-file-list="false"
              :with-credentials="true"
              :multiple="true"
              :http-request="addFile"
              :accept="fileAccept"
          >
            <div class="op-item">
              <Icon iconName="file" :width="60"></Icon>
              <div>上传文件</div>
            </div>
          </el-upload>
          <div class="op-item" v-if="category === 'all'" @click="newFolder">
            <Icon iconName="folder" :width="60"></Icon>
            <div>新建目录</div>
          </div>
        </div>
      </div>
    </div>
    <!--预览-->
    <Preview ref="previewRef"></Preview>
    <!--移动-->
    <FolderSelect
        ref="folderSelectRef"
        @folderSelect="moveFolderDone"
    ></FolderSelect>
    <!--分享-->
    <FileShare ref="shareRef"></FileShare>
  </div>
</template>

<script setup>
import CategoryInfo from "@/js/CategoryInfo.js";
import FileShare from "./ShareFile.vue";
import {ref, reactive, getCurrentInstance, nextTick, computed, onMounted} from "vue";
import {useRouter, useRoute} from "vue-router";
import {UploadFilled} from '@element-plus/icons-vue'

const {proxy} = getCurrentInstance();
const router = useRouter();
const route = useRoute();
const emit = defineEmits(["addFile"]);
//import download from "../../Func/download";
import exportData from "../../Func/download";
import Icon from "../../components/Icon.vue";
import axios from "axios";
//import {dialogEmits as upload} from "../../../.vite/deps/element-plus";

const downloadFile = async () => {
  console.log("执行")
  await exportData()
}
/*** 用户导入参数 */
const upload = reactive({
  // 是否显示弹出层（用户导入）
  open: false,
  // 弹出层标题（批量邀约）
  title: "",
  // 是否禁用上传
  isUploading: false,
  // 是否更新已经存在的用户数据
  //updateSupport: 0,
  // 设置上传的请求头部
  headers: {"token": localStorage.getItem("token")},
  // 上传的地址
  url: "https://localhost/user/file/large"
});

/** 下载模板操作 */
function importTemplate() {
  proxy.download("system/user/importTemplate", {}, `user_template_${new Date().getTime()}.xlsx`);
}

/**文件上传中处理 */
const handleFileUploadProgress = (event, file, fileList) => {
  upload.isUploading = true;
};
/** 文件上传成功处理 */
const handleFileSuccess = (response, file, fileList) => {
  upload.open = false;
  upload.isUploading = false;
  proxy.$refs["uploadRef"].handleRemove(file);
  proxy.$alert("<div style='overflow: auto;overflow-x: hidden;max-height: 70vh;padding: 10px 20px 0;'>" + response.msg + "</div>", "导入结果", {dangerouslyUseHTMLString: true});
  getList();
};

/** 提交上传文件 */
function submitFileForm() {
  proxy.$refs["uploadRef"].submit();
}


const addFile = async (fileData) => {
  emit("addFile", {file: fileData.file, filePid: currentFolder.value.fileId});
};


//添加文件回调
const reload = () => {
  showLoading.value = false;
  loadDataList();
};
defineExpose({
  reload,
});

//当前文件夹
const currentFolder = ref({fileId: 0});

const api = {
  loadDataList: "/file/loadDataList",
  rename: "/file/rename",
  newFoloder: "/file/newFoloder",
  getFolderInfo: "/file/getFolderInfo",
  delFile: "/file/delFile",
  changeFileFolder: "/file/changeFileFolder",
  createDownloadUrl: "/file/createDownloadUrl",
  //download: "/api/file/download",
};

const fileAccept = computed(() => {
  const categoryItem = CategoryInfo[category.value];
  return categoryItem ? categoryItem.accept : "*";
});

//列表
const columns = [
  {
    label: "文件名",
    prop: "fileName",
    scopedSlots: "fileName",
  },
  {
    label: "修改时间",
    prop: "lastUpdateTime",
    width: 200,
  },
  {
    label: "大小",
    prop: "fileSize",
    scopedSlots: "fileSize",
    width: 200,
  },
];
//搜索
const search = () => {
  showLoading.value = true;
  loadDataList();
};
//列表 //文件列表来源
const tableData = ref({});
const tableOptions = {
  extHeight: 50,
  selectType: "checkbox",
};

const fileNameFuzzy = ref();
const showLoading = ref(true);
const category = ref();

const getFile = () => {
  console.log("执行");
  axios.get('https://localhost/user/dir/all/earliest/1', {
    headers: {
      'Content-Type': 'application/json',
      'token': localStorage.getItem("token"),
    }
  })
      .then(response => {
        // 处理响应
        console.log(response)


      })
      .catch(error => {
        // 处理错误
        console.error(error);
      });

}
onMounted(getFile);
//展示文件菜单
const loadDataList = async () => {
  let params = {
    pageNo: tableData.value.pageNo,
    pageSize: tableData.value.pageSize,
    fileNameFuzzy: fileNameFuzzy.value,
    category: category.value,
    filePid: currentFolder.value.fileId,
  };
  if (params.category !== "all") {
    delete params.filePid;
  }
  // let result = await proxy.Request({
  //   url: api.loadDataList,
  //   showLoading: showLoading,
  //   params,
  // });
  let result = null;
  if (!result) {
    return;
  }
  tableData.value = result.data;
  editing.value = false;
};

//展示操作按钮
const showOp = (row) => {
  tableData.value.list.forEach((element) => {
    element.showOp = false;
  });
  row.showOp = true;
};

const cancelShowOp = (row) => {
  row.showOp = false;
};

//编辑行
const editing = ref(false);
const editNameRef = ref();
//新建文件夹
const newFolder = () => {
  if (editing.value) {
    return;
  }
  tableData.value.list.forEach((element) => {
    element.showEdit = false;
  });
  editing.value = true;
  tableData.value.list.unshift({
    showEdit: true,
    fileType: 0,
    fileId: "",
    filePid: currentFolder.value.fileId,
  });
  nextTick(() => {
    editNameRef.value.focus();
  });
};

const cancelNameEdit = (index) => {
  const fileData = tableData.value.list[index];
  if (fileData.fileId) {
    fileData.showEdit = false;
  } else {
    tableData.value.list.splice(index, 1);
  }
  editing.value = false;
};

const saveNameEdit = async (index) => {
  const {fileId, filePid, fileNameReal} = tableData.value.list[index];
  if (fileNameReal == "" || fileNameReal.indexOf("/") != -1) {
    proxy.Message.warning("文件名不能为空且不能含有斜杠");
    return;
  }
  let url = api.rename;
  if (fileId == "") {
    url = api.newFoloder;
  }
  // let result = await proxy.Request({
  //   url: url,
  //   params: {
  //     fileId,
  //     filePid: filePid,
  //     fileName: fileNameReal,
  //   },
  // });
  let result = null
  if (!result) {
    return;
  }
  tableData.value.list[index] = result.data;
  editing.value = false;
};

//编辑文件名
const editFileName = (index) => {
  if (tableData.value.list[0].fileId == "") {
    tableData.value.list.splice(0, 1);
    index = index - 1;
  }
  tableData.value.list.forEach((element) => {
    element.showEdit = false;
  });
  let cureentData = tableData.value.list[index];
  cureentData.showEdit = true;

  //编辑文件
  if (cureentData.folderType == 0) {
    cureentData.fileNameReal = cureentData.fileName.substring(
        0,
        cureentData.fileName.indexOf(".")
    );
    cureentData.fileSuffix = cureentData.fileName.substring(
        cureentData.fileName.indexOf(".")
    );
  } else {
    cureentData.fileNameReal = cureentData.fileName;
    cureentData.fileSuffix = "";
  }
  editing.value = true;
  nextTick(() => {
    editNameRef.value.focus();
  });
};

//多选 批量选择
const selectFileIdList = ref([]);
const selectFileList = ref([]);
const rowSelected = (rows) => {
  selectFileList.value = rows;
  selectFileIdList.value = [];
  rows.forEach((item) => {
    selectFileIdList.value.push(item.fileId);
  });
};

//删除文件
const delFile = (row) => {
  proxy.Confirm(
      `你确定要删除【${row.fileName}】吗？删除的文件可在10天内通过回收站还原`,
      async () => {
        let result = null;
        //     await proxy.Request({
        //   url: api.delFile,
        //   params: {
        //     fileIds: row.fileId,
        //   },
        // });
        if (!result) {
          return;
        }
        loadDataList();
      }
  );
};
//批量删除
const delFileBatch = () => {
  if (selectFileIdList.value.length == 0) {
    return;
  }
  proxy.Confirm(
      `你确定要删除这些文件吗？删除的文件可在10天内通过回收站还原`,
      async () => {
        let result = await proxy.Request({
          url: api.delFile,
          params: {
            fileIds: selectFileIdList.value.join(","),
          },
        });
        if (!result) {
          return;
        }
        loadDataList();
      }
  );
};

//移动目录
const folderSelectRef = ref();
const currentMoveFile = ref({});
const moveFolder = (data) => {
  currentMoveFile.value = data;
  folderSelectRef.value.showFolderDialog(data.fileId);
};

//批量移动
const moveFolderBatch = () => {
  currentMoveFile.value = {};
  //批量移动如果选择的是文件夹，那么要讲文件夹也过滤
  const excludeFileIdList = [currentFolder.value.fileId];
  selectFileList.value.forEach((item) => {
    if (item.folderType == 1) {
      excludeFileIdList.push(item.fileId);
    }
  });
  folderSelectRef.value.showFolderDialog(excludeFileIdList.join(","));
};

const moveFolderDone = async (folderId) => {
  if (
      currentMoveFile.value.filePid === folderId ||
      currentFolder.value.fileId == folderId
  ) {
    proxy.Message.warning("文件正在当前目录，无需移动");
    return;
  }
  let filedIdsArray = [];
  if (currentMoveFile.value.fileId) {
    filedIdsArray.push(currentMoveFile.value.fileId);
  } else {
    filedIdsArray = filedIdsArray.concat(selectFileIdList.value);
  }
  let result = await proxy.Request({
    url: api.changeFileFolder,
    params: {
      fileIds: filedIdsArray.join(","),
      filePid: folderId,
    },
  });
  if (!result) {
    return;
  }
  folderSelectRef.value.close();
  loadDataList();
};

const previewRef = ref();
const navigationRef = ref();
const preview = (data) => {
  if (data.folderType == 1) {
    //openFolder(data);
    navigationRef.value.openFolder(data);
    return;
  }
  if (data.status != 2) {
    proxy.Message.warning("文件正在转码中，无法预览");
    return;
  }
  previewRef.value.showPreview(data, 0);
};

//目录
const navChange = (data) => {
  const {curFolder, categoryId} = data;
  currentFolder.value = curFolder;
  showLoading.value = true;
  category.value = categoryId;
  loadDataList();
};

// //下载文件
// const download = async (row) => {
//   let result = await proxy.Request({
//     url: api.createDownloadUrl + "/" + row.fileId,
//   });
//   if (!result) {
//     return;
//   }
//   window.location.href = api.download + "/" + result.data;
// };

//分享
const shareRef = ref();
const share = (row) => {
  shareRef.value.show(row);
};
</script>

<style lang="scss" scoped>
@import "../../assets/file.list.scss";
</style>