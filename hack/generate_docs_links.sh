DOCS_PATH="docs"

lines=$(ls $DOCS_PATH | wc -w)
lines=$(expr $lines + 2)
s_pos=$(grep -n '### Commands' README.md | cut -d':' -f1)

if [ ! -z $s_pos ]
then
  e_pos=$(expr $s_pos + $lines)
  sed -i ${s_pos},${e_pos}d README.md
fi

echo "### Commands" >> README.md
echo "" >> README.md
for file in $(ls $DOCS_PATH)
do
  link=$(echo "$file" | sed 's/\.md//' | sed 's/_/ /g')
  echo "- [$link](./$DOCS_PATH/$file)" >> README.md
done

