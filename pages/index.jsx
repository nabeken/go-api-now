import { useState } from 'react';
import Head from 'next/head';
import useSWR from 'swr';

const fetcher = (url) => fetch(url).then((res) => res.json());

export default function Homepage() {
  const names = ['Learn', 'Develop', 'Preview', 'Ship'];

  return (
    <>
      <Head>
        <title>go-api-now</title>
      </Head>

      {names.map((name) => (
        <Header key={ name } title={ `${name}.` }  />
      ))}

      { /*
      <Header title="Learn." />
      <Header title="Learn. Develop." />
      <Header title="Learn. Develop. Preview." />
      <Header title="Learn. Develop. Preview. Ship. 🚀" />
      */}

      <p>Next.js</p>

      <LikeButton names={names} />
      <LikeButton names={names} />

    <h1>Ping</h1>
    <pre><Ping /></pre>
    </>
  );
}

function Ping() {
  const { data, error, isLoading } = useSWR('/api/ping', fetcher);

  if (error) return <>failed to load</>;
  if (isLoading) return <>loading...</>;

  return <>hello, the API version is {data.version} at {data.now}!</>;
}

function LikeButton({ names }) {
  const [idx, setIdx] = useState(0);
  const hasNext = idx < names.length - 1;

  function handleClick() {
    if (hasNext) {
      setIdx(idx+1);
    } else {
      setIdx(0);
    }
  }

  return <button onClick={handleClick}>Like ({names[idx]})</button>;
}

function Header({ title }) {
  return <h1>{title}</h1>;
}
