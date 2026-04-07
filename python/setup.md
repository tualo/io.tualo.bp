brew install pyenv
pyenv install 3.9.10
pyenv global 3.9.10

PATH=~/.pyenv/versions/3.9.10/bin/:$PATH

# PATH=~/.pyenv/versions/3.9.10/bin/:/opt/homebrew/opt/mysql-client/bin:~/go/bin/:/Applications/Visual\ Studio\ Code.app/Contents/Resources/app/bin/:/Applications/Visual\ Studio\ Code.app/Contents/Resources/app/bin/:/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/usr/local/MacGPG2/bin:/usr/local/go/bin:/opt/homebrew/opt/mysql-client/bin:~/go/bin/:/Applications/Visual\ Studio\ Code.app/Contents/Resources/app/bin/


python3 -m pip install tensorflow-macos
# python3 -m pip install tensorflow-metal
python3 -m pip install matplotlib

mkdir -p data/O
mkdir -p data/X
mkdir -p data/Q

go run main.go "thomashoffmann:@tcp(127.0.0.1:3306)/bwnuernberg" "17680"

python python/train3.py 

python python/server.py 