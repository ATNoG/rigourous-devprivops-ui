rm -rf /tmp/.devprivops 
mkdir /tmp/.devprivops 
cp -r .devprivops/ /tmp/.devprivops/test 

alias git='git -C /tmp/.devprivops/test'

git init
git config --local user.name "Admin"
git config --local user.email "admin@corp.co"
git config --local receive.denyCurrentBranch updateInstead

git add .
git commit -m "First commit"

