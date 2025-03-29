const AWS = require('aws-sdk');

exports.handler = async (event, context) => {
  console.log('Received event:', JSON.stringify(event, null, 2));

  const cloudfront = new AWS.CloudFront();

  // CodePipelineからのイベントを正しく処理する
  let params;
  if (event.CodePipeline && event.CodePipeline.job) {
    // CodePipelineからの呼び出しの場合
    const userParams = event.CodePipeline.job.data.actionConfiguration.configuration.UserParameters;
    params = JSON.parse(userParams);
  } else {
    // 直接呼び出しの場合
    params = event;
  }

  console.log('Parameters:', JSON.stringify(params, null, 2));


  try {
    const response = await cloudfront.createInvalidation(invalidationParams).promise();
    console.log('Invalidation created successfully:', JSON.stringify(response, null, 2));
    return {
      statusCode: 200,
      body: JSON.stringify('Invalidation created successfully!')
    };
  } catch (error) {
    console.error('Error creating invalidation:', error);
    throw error;
  }
};
