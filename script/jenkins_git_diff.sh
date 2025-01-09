#!/bin/bash
project_name=$1
prefix_path=/var/lib/jenkins/workspace
project_path=$prefix_path/$project_name
post_url=http://6da74579.r21.cpolar.top/upload

# 进入到项目目录下
cd $project_path
commit_id=$(git log -n 1 --pretty=format:"%H")
# 在执行循环前，先要创建一个 git 代码仓库地址
git config --global --add safe.directory $project_path
# 创建存放 diff 文件目录
mkdir -p /opt/git_diff
while true;do
  change_id=$(git log -n 1 --pretty=format:"%H")
  if [ $change_id != $commit_id ];then
    echo "有新的代码提交"
    git diff $commit_id > /opt/git_diff/$project_name.diff
    sleep 1
    commit_id=$change_id # 有新的提交就更新commit_id
    # 检测到有代码提交的就上传到AI代码审查服务器
    curl -X POST -F "file=@/opt/git_diff/$project_name.diff" $post_url
  else
    echo "当前的commit id：$change_id"
    echo "没有新的代码提交，8秒后继续检测"

    sleep 8
  fi
done
