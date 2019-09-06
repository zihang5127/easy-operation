function loadCodeMirrorEditor(id) {
    window.CodeMirrorEditor = CodeMirror.fromTextArea(document.getElementById(id),{
        mode: "shell",    //实现groovy代码高亮
        lineNumbers: true,	//显示行号
        theme: "dracula",	//设置主题
        lineWrapping: true,	//代码折叠
        foldGutter: true,
        gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
        matchBrackets: true,	//括号匹配
    });

    window.CodeMirrorEditor.on("change",function () {
        var content = window.CodeMirrorEditor.getValue();
        $("#" + id).val(content);
    });
}

(function ($) {

    $('[data-toggle="tooltip"]').tooltip();


    $.fn.open = function (title, values) {
        this.find(".modal-title").text(title);
        if(typeof values === "object"){
            for (var i in values){

            }
        }else{
            this.find(".modal-body").html(values)
        }
        return this;
    };

    var projectModal = $("#projectModal");
    var projectCache = projectModal.find(".modal-body").html();
    var serverModal = $("#serverModal");
    var serverCache = serverModal.find(".modal-body").html();

    $("#addProjectBtn").on("click",function () {
        projectModal.open("New Project",projectCache).modal("show");
    });
    projectModal.on("show.bs.modal",function () {
        $('[data-toggle="tooltip"]').tooltip();
    });
    projectModal.on("shown.bs.modal",function () {
        if(!window.CodeMirrorEditor){
            window.loadCodeMirrorEditor("shellScript");
        }
    });

    $("#addServerBtn").on("click",function () {
        serverModal.open("New Server",serverCache).modal("show");
    });
    $("#serverForm").ajaxForm({
        beforeSubmit : function () {
            var serverNameEle = $("#serverName");

            if($.trim(serverNameEle.val()) === ""){
                serverNameEle.closest(".form-group").addClass("has-error");
                return false;
            }

            var serverIp = $("#serverIp");

            if($.trim(serverIp.val()) === ""){
                serverIp.closest(".form-group").addClass("has-error");
                return false;
            }
            var serverPort = $("#serverPort");
            if($.trim(serverPort.val()) === ""){
                serverPort.closest(".form-group").addClass("has-error");
                return false;
            }

            var serverKey = $("#serverKey");
            if($.trim(serverKey.val()) === ""){
                serverKey.closest(".form-group").addClass("has-error");
                return false;
            }

            $("#saveServerBtn").button("load");
        } ,
        success :function (res) {
            if (res.errcode === 0){
                if (serverModal.length > 0) {

                    $("#serverTable>tbody").prepend(res.view);
                    serverModal.modal("hide");
                }else{
                    $("#errorMessage").css("color","green").text("success");
                }
            }else {
                $("#errorMessage").css("color","red").text(res.message);
            }
        },complete : function () {
            $("#saveServerBtn").button("reset");
        }
    });

    $("#projectForm").ajaxForm({
        dataType :"json",
        beforeSubmit :function () {
            var isValid = true;
            var repositoryNameEle = $("#repo_name");
            if($.trim(repositoryNameEle.val()) == ""){
                repositoryNameEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }
            var repositoryBranchEle = $("#branch_name");
            if($.trim(repositoryBranchEle.val()) === ""){
                repositoryBranchEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }

            var repositoryTypeEle = $("#repositoryType");
            if($.trim(repositoryTypeEle.val()) === ""){
                repositoryTypeEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }
            var shellScriptEle = $("#shellScript");
            if($.trim(shellScriptEle.val()) === ""){
                shellScriptEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }

            if(!isValid){
                return false;
            }
            $("#saveProjectBtn").button("load");
        },success : function (res) {
            if (res.errcode === 0){
                if (projectModal.length > 0) {
                    $("#projectTable>tbody").prepend(res.view);
                    projectModal.modal("hide");
                }else{
                    $("#errorMessage").css("color","green").text("success");
                }
            }else {
                $("#errorMessage").css("color","red").text(res.message);
            }
        },complete : function () {
            $("#saveProjectBtn").button("reset");
        },
        error :function () {
            $("#errorMessage").css("color","red").text("Server error. Please refresh and try again.");
        }
    });


})(jQuery);