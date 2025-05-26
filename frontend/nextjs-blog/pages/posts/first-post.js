import Link from "next/link";
import Head from "next/head";
// import Script from 'next/script';
import Layout from "../../components/layout";
import Alert  from "../../components/alert";    

export default function FirstPost() {
  return (
    <Layout home={false}>
      <Head>
        <title>First Post</title>

        {/* <script src="https://connect.facebook.net/en_US/sdk.js" />  */}

      </Head>

      {/* <Script
        src="https://connect.facebook.net/en_US/sdk.js"
        strategy="lazyOnload"
        onLoad={() =>
          console.log(`腳本正確地載入，window.FB 已被定義`)
        }
      /> */}

      <h1>First Post</h1>
      <h2>
        <Link href="/">Back to home</Link>
      </h2>
      <Alert type="success">This is a success alert</Alert>
    </Layout>
  );
}
