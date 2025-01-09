#!/bin/bash
project_id=$1
branch=$2
save_path=$3
post_url=http://6da74579.r21.cpolar.top/upload

# 脚本使用示例：
# ./gitlab_git_diff.sh "项目ID" "branch名" "保存git-diff文件的路径"

# 进入项目ID的sha256编码
sha256_ID=$(echo -n $project_id | sha256sum | awk '{print $1}')
dir_1=$(echo $sha256_ID | cut -c 1-2)
dir_2=$(echo $sha256_ID | cut -c 3-4)
gitlab_repo=/var/opt/gitlab/git-data/repositories/@hashed/$dir_1/$dir_2/$sha256_ID.git

# 进入到对应的代码仓库
cd "$gitlab_repo" || { echo "无法进入目录: $gitlab_repo"; exit 1; }
# 添加到 git 仓库，方便后面进行 git 操作
git config --global --add safe.directory "$gitlab_repo"


# 判断是否存在 branch.git 目录
if [ -d "$gitlab_repo/$branch.git" ];then 
    echo "目标目录 $gitlab_repo/$branch.git 已存在，现已删除..."
    rm -rf "$gitlab_repo/$branch.git"
fi
# clone 对应的分支到当前目录（创建工作树）
git clone -b "$branch" "$gitlab_repo" "$gitlab_repo/$branch.git"
sleep 2
# 进入到对应的分支里面
cd $gitlab_repo/$branch.git
# 添加到 git 仓库，方便后面进行 git 操作
git config --global --add safe.directory "$gitlab_repo/$branch.git"
# 获取当前的 Commit ID 
commit_id=$(git log -n 1 --pretty=format:"%H")
# 创建存放 diff 文件目录
mkdir -p /opt/git_diff

while true;do
    # 进入到对应的分支里面
    cd $gitlab_repo/$branch.git
    change_id=$(git log -n 1 --pretty=format:"%H")
    sleep 2
    if [ $change_id != $commit_id ];then
    echo "有新的代码提交"
    git diff $commit_id > $save_path/$project_id.diff
    cd $gitlab_repo
    sleep 1
    commit_id=$change_id # 有新的提交就更新commit_id
    # 追加上一次的commit id 到 diff 文件方便读取
    echo "commitID:$commit_id" >> $save_path/$project_id.diff
    # 然后再把项目的 ID 传进去
    echo "projectID:$project_id" >> $save_path/$project_id.diff
    # 检测到有代码提交的就上传到AI代码审查服务器
    curl -X POST -F "file=@$save_path/$project_id.diff" $post_url
else
    cd $gitlab_repo
    echo "当前的commit id：$change_id"
    echo "没有新的代码提交，20秒后继续检测"
    if [ -d "$gitlab_repo/$branch.git" ];then 
        echo "目标目录 $gitlab_repo/$branch.git 已存在，现已删除..."
        rm -rf "$gitlab_repo/$branch.git"
    fi
    git clone -b "$branch" "$gitlab_repo" "$gitlab_repo/$branch.git"
    sleep 20
    fi
done