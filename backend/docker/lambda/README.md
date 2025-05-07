build push方法

backendのディレクトリで作業した。コマンドはうろ覚えだが、buildしてタグ付けしてpushする流れでやった。

docker build -t tasmap:latest -f docker/lambda/Dockerfile .

docker tag tasmap:latest 420302062688.dkr.ecr.ap-northeast-1.amazonaws.com/tasmap-lambda:latest

awsのecrにログイン

docker push 420302062688.dkr.ecr.ap-northeast-1.amazonaws.com/tasmap-lambda:latest

