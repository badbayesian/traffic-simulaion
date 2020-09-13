#!/bin/bash
#SBATCH --job-name=proj3                      # Job name
#SBATCH --mail-type=ALL                     # Mail events (NONE, BEGIN, END, FAIL, ALL)
#SBATCH --mail-user=badbayesian@gmail.com   # Where to send mail
#SBATCH --ntasks=8                          # Run on a single CPU
#SBATCH --partition=fast
#SBATCH --mem=8gb                           # Job memory request
#SBATCH --time=04:30:00                     # Time limit hrs:min:sec

module load golang/1.14.1

rm ../report/report.txt

echo runtime model thread car n  >> ../report/report.txt

threads=(1 2 4 6 8)
models=("city" "complete" "random")
experiments=20
rounds=100
cars=(10 100 1000)
size=(10 100 1000)

cd ../main
for model in "${models[@]}"; do
	for thread in "${threads[@]}"; do
		for car in "${cars[@]}"; do
			for n in "${size[@]}"; do
				start=`date +%s.%N`
				go run simulation.go --p=$thread --n=$experiments --r=$rounds --c=$car --s=1 --size=$n --t=$model --loc=wut
			end=`date +%s.%N`
			runtime=$( echo "$end - $start" | bc -l )
			echo $runtime $model $thread $car $n >> ../report/report.txt
		done
	done
done
done
