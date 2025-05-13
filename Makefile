run-node:
	sudo docker run -p 9944:9944 -p 9615:9615 parity/polkadot:v1.16.2 --name "my-node" --rpc-external --prometheus-external